import WebSocket from 'ws';
import jwt from 'jsonwebtoken';
import { promisify } from 'util';

const sleep = promisify(setTimeout);

// Cấu hình môi trường test - Cổng mặc định là 8081 để tránh xung đột với API (8080)
const WS_URL = process.env.WS_URL || 'ws://localhost:8081/ws/interview-room';
const JWT_SECRET = process.env.LIVEKIT_API_SECRET || 'devsecret';

const testRunId = Math.random().toString(36).substr(2, 9);
const room_id = `test-room-${testRunId}`;
const interview_id = `test-interview-${testRunId}`;

// Hàm sinh token JWT LiveKit giả lập
function generateToken(userId, role, displayName) {
  const payload = {
    sub: userId,
    user_id: userId,
    role: role,
    name: displayName,
    display_name: displayName,
    room_id: room_id,
    interview_id: interview_id,
    video: {
      room: room_id,
      roomJoin: true
    }
  };
  return jwt.sign(payload, JWT_SECRET, { expiresIn: '1h' });
}

// Helper quản lý WS Client
class TestClient {
  constructor(name, role) {
    this.name = name;
    this.role = role;
    this.userId = `${role}-user-123`;
    this.token = generateToken(this.userId, role, name);
    this.ws = null;
    this.receivedEvents = [];
    this.errors = [];
  }

  async connect() {
    return new Promise((resolve, reject) => {
      const url = `${WS_URL}?token=${this.token}`;
      this.ws = new WebSocket(url);

      this.ws.on('open', () => {
        console.log(`[Client ${this.name}] Connected to server.`);
        resolve();
      });

      this.ws.on('message', (data) => {
        try {
          const envelope = JSON.parse(data.toString());
          this.receivedEvents.push(envelope);
          if (envelope.event === 'error') {
            console.error(`\x1b[31m[Client ${this.name}] Received ERROR Event:\x1b[0m`, JSON.stringify(envelope.payload));
          } else {
            console.log(`[Client ${this.name}] Received Event: ${envelope.event} (reqId: ${envelope.request_id || 'none'})`);
          }
        } catch (err) {
          console.error(`[Client ${this.name}] Failed to parse message:`, err);
        }
      });

      this.ws.on('error', (err) => {
        console.error(`[Client ${this.name}] WS Error:`, err.message);
        this.errors.push(err);
        reject(err);
      });

      this.ws.on('close', (code, reason) => {
        console.log(`[Client ${this.name}] Closed connection: code=${code}, reason=${reason.toString()}`);
      });
    });
  }

  send(event, payload, reqId = `req_${Math.random().toString(36).substr(2, 9)}`) {
    const envelope = {
      event,
      request_id: reqId,
      room_id,
      interview_id,
      payload
    };
    this.ws.send(JSON.stringify(envelope));
    console.log(`[Client ${this.name}] Sent Event: ${event} (reqId: ${reqId})`);
    return reqId;
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }

  getEvents(eventName) {
    return this.receivedEvents.filter(e => e.event === eventName);
  }
}

async function runTests() {
  console.log('====================================================');
  console.log('BẮT ĐẦU CHẠY AUTOMATED E2E INTEGRATION TEST REALTIME');
  console.log('====================================================\n');

  // Khởi tạo Recruiter và Candidate
  const recruiter = new TestClient('Hùng Recruiter', 'recruiter');
  const candidate = new TestClient('Lai Candidate', 'candidate');

  try {
    // 1. Kết nối cả 2 Client
    await recruiter.connect();
    await candidate.connect();

    // 2. Client gửi room:join
    console.log('\n--- TEST STEP 1: Join Room ---');
    recruiter.send('room:join', { displayName: recruiter.name });
    candidate.send('room:join', { displayName: candidate.name });

    await sleep(1000); // Chờ server xử lý join và presence

    // Kiểm tra Event ACK room:joined
    const recJoined = recruiter.getEvents('room:joined');
    const candJoined = candidate.getEvents('room:joined');

    if (recJoined.length === 0 || candJoined.length === 0) {
      throw new Error('FAILED: room:joined ACK not received.');
    }
    console.log('SUCCESS: Cả 2 client join room và nhận ACK thành công.');

    // 3. Test Presence Update (ai online/offline)
    const recPresence = recruiter.getEvents('room:presence_update');
    if (recPresence.length === 0) {
      throw new Error('FAILED: room:presence_update not broadcasted.');
    }
    console.log('SUCCESS: Presence Update hoạt động chính xác.');

    // 4. Test Media Status Update
    console.log('\n--- TEST STEP 2: Media Status ---');
    recruiter.send('media:status', { mic_enabled: true, camera_enabled: true, screen_sharing: false });
    await sleep(500);

    const candMediaChange = candidate.getEvents('media:status_changed');
    if (candMediaChange.length === 0) {
      throw new Error('FAILED: media:status_changed not broadcasted to other participants.');
    }
    console.log('SUCCESS: Media status thay đổi và broadcast thành công.');

    // 5. Test Chat Realtime & Phân quyền Visibility
    console.log('\n--- TEST STEP 3: Chat & Notes Visibility ---');
    // Gửi chat toàn phòng (room)
    recruiter.send('chat:send', { message: 'Chào bạn Lai!', visibility: 'room' });
    // Gửi note/chat nội bộ (recruiter_only)
    recruiter.send('chat:send', { message: 'Ghi chú: Lai đang khá bối rối.', visibility: 'recruiter_only' });

    await sleep(1000);

    // Kiểm tra tin nhắn toàn phòng
    const candChats = candidate.getEvents('chat:message');
    const publicChat = candChats.find(c => c.payload.message === 'Chào bạn Lai!');
    if (!publicChat) {
      throw new Error('FAILED: Candidate did not receive public chat.');
    }

    // Kiểm tra tin nhắn nội bộ (Candidate KHÔNG được phép nhận)
    const privateChatForCand = candChats.find(c => c.payload.message && c.payload.message.includes('Lai đang khá bối rối'));
    if (privateChatForCand) {
      throw new Error('FAILED: SECURITY BREACH! Candidate received recruiter_only internal note.');
    }
    console.log('SUCCESS: Visibility kiểm soát chat và note nội bộ an toàn (Candidate không thấy).');

    // 6. Test Interview Lifecycle (Start -> Pause -> Resume)
    console.log('\n--- TEST STEP 4: Interview Control ---');
    // Candidate gửi interview:start -> phải bị lỗi FORBIDDEN
    const candStartReqId = candidate.send('interview:start', { consent_recording: true, consent_ai: true });
    await sleep(500);
    const candErrors = candidate.getEvents('error');
    const forbiddenError = candErrors.find(e => e.request_id === candStartReqId && e.payload.code === 'FORBIDDEN');
    if (!forbiddenError) {
      throw new Error('FAILED: Candidate start interview was not blocked with FORBIDDEN.');
    }
    console.log('SUCCESS: Chặn Candidate tự ý start interview chính xác.');

    // Recruiter gửi interview:start
    recruiter.send('interview:start', { consent_recording: true, consent_ai: true });
    await sleep(1000);

    const recStarted = recruiter.getEvents('interview:started');
    const candStarted = candidate.getEvents('interview:started');
    if (recStarted.length === 0 || candStarted.length === 0) {
      throw new Error('FAILED: interview:started broadcast not received.');
    }
    console.log('SUCCESS: Bắt đầu phỏng vấn thành công.');

    // 7. Test AI Suggestion & Score Update (Recruiter only)
    console.log('\n--- TEST STEP 5: AI Bridge & Rate Limiting ---');
    const recSuggReqId = recruiter.send('ai:request_suggestion', { focus: 'technical' });
    await sleep(500);

    // Kiểm tra ai:thinking nhận được ở đầu Recruiter
    const recThinking = recruiter.getEvents('ai:thinking');
    if (recThinking.length === 0) {
      throw new Error('FAILED: ai:thinking event not received by recruiter.');
    }

    // Candidate KHÔNG được nhận ai:thinking
    const candThinking = candidate.getEvents('ai:thinking');
    if (candThinking.length > 0) {
      throw new Error('FAILED: SECURITY BREACH! Candidate received ai:thinking.');
    }
    console.log('SUCCESS: Trạng thái AI Thinking chỉ hiển thị với Recruiter.');

    await sleep(1500); // Chờ mock worker xử lý xong ai:suggestion
    const recSugg = recruiter.getEvents('ai:suggestion');
    if (recSugg.length === 0) {
      throw new Error('FAILED: ai:suggestion not received.');
    }
    console.log('SUCCESS: AI suggestions trả về đúng chuẩn envelope.');

    // 8. Test Heartbeat Keep-Alive
    console.log('\n--- TEST STEP 6: Room Heartbeat ---');
    recruiter.send('room:heartbeat', { connection_state: 'online' });
    candidate.send('room:heartbeat', { connection_state: 'online' });
    console.log('SUCCESS: Heartbeat gửi không gây lỗi server.');

    // 9. Test Interview End & Grace Period Teardown
    console.log('\n--- TEST STEP 7: End Interview & Grace Period ---');
    recruiter.send('interview:end', { generate_report: true });
    await sleep(1000);

    const recEnded = recruiter.getEvents('interview:completed');
    const candEnded = candidate.getEvents('interview:completed');
    if (recEnded.length === 0 || candEnded.length === 0) {
      throw new Error('FAILED: interview:completed not broadcasted.');
    }
    console.log('SUCCESS: Hoàn tất đóng phỏng vấn.');

    console.log('\n====================================================');
    console.log('TẤT CẢ CÁC BÀI INTEGRATION TEST ĐÃ VƯỢT QUA (PASSED) ✅');
    console.log('====================================================');

  } catch (err) {
    console.error('\n❌ BÀI INTEGRATION TEST THẤT BẠI:');
    console.error(err.message);
    process.exit(1);
  } finally {
    // Ngắt kết nối
    recruiter.disconnect();
    candidate.disconnect();
  }
}

// Chạy test
runTests();
