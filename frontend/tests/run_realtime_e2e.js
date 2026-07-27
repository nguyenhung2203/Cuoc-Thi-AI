import { spawn, execSync } from 'child_process';
import path from 'path';
import fs from 'fs';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Đường dẫn tương đối đến backend và script test
const BACKEND_DIR = path.resolve(__dirname, '../../backend');
const FRONTEND_DIR = path.resolve(__dirname, '../');
const TEST_SCRIPT = path.resolve(__dirname, './realtime_integration_test.js');
const BINARY_NAME = 'realtime_server_temp.exe';
const BINARY_PATH = path.join(BACKEND_DIR, BINARY_NAME);
const PORT = '8081'; // Sử dụng cổng 8081 cho Realtime Gateway

async function main() {
  console.log('====================================================');
  console.log('KHỞI ĐỘNG HỆ THỐNG KIỂM THỬ TỰ ĐỘNG REALTIME');
  console.log('====================================================\n');

  // 1. Biên dịch Go Backend Realtime Server trước để tránh tiến trình mồ côi (dangling process)
  console.log('[BE] Đang biên dịch Go Realtime server...');
  try {
    if (fs.existsSync(BINARY_PATH)) {
      fs.unlinkSync(BINARY_PATH);
    }
    execSync(`go build -o ${BINARY_NAME} ./cmd/realtime/main.go`, { cwd: BACKEND_DIR, stdio: 'inherit' });
    console.log('[BE] Biên dịch thành công.');
  } catch (err) {
    console.error('❌ Lỗi biên dịch Go backend:', err.message);
    process.exit(1);
  }

  console.log(`[BE] Đang khởi động Go Realtime server binary trên cổng ${PORT}...`);
  const server = spawn(BINARY_PATH, [], {
    cwd: BACKEND_DIR,
    env: { 
      ...process.env, 
      LIVEKIT_API_SECRET: 'devsecret',
      REALTIME_PORT: PORT
    }
  });

  let serverStarted = false;
  let serverExited = false;
  let serverExitCode = null;

  server.on('close', (code) => {
    serverExited = true;
    serverExitCode = code;
  });
  
  // Lắng nghe stdout/stderr của Backend
  server.stdout.on('data', (data) => {
    const output = data.toString();
    const lines = output.trim().split('\n');
    lines.forEach(line => {
      console.log(`\x1b[36m[BE LOG]\x1b[0m ${line}`);
    });

    if (output.includes('listening on') || output.includes('Starting AI Interview')) {
      serverStarted = true;
    }
  });

  server.stderr.on('data', (data) => {
    const output = data.toString();
    const lines = output.trim().split('\n');
    lines.forEach(line => {
      console.error(`\x1b[31m[BE ERR]\x1b[0m ${line}`);
    });
    if (output.includes('listening on')) {
      serverStarted = true;
    }
  });

  // Chờ server khởi động thành công (tối đa 15s)
  for (let i = 0; i < 30; i++) {
    if (serverStarted) break;
    if (serverExited) {
      console.error(`❌ Lỗi: Backend server đã thoát đột ngột với mã ${serverExitCode} trước khi kịp lắng nghe. Hãy đảm bảo không có tiến trình nào khác đang chiếm dụng cổng :${PORT}.`);
      process.exit(1);
    }
    await new Promise(r => setTimeout(r, 500));
  }

  if (!serverStarted) {
    console.error(`❌ Lỗi: Backend server không khởi động thành công sau 15 giây.`);
    server.kill();
    process.exit(1);
  }

  console.log('\n[BE] Server đã lắng nghe. Đang chạy test client...');

  // 2. Chạy test script client
  const testProcess = spawn('node', [TEST_SCRIPT], {
    cwd: FRONTEND_DIR,
    env: {
      ...process.env,
      WS_URL: `ws://localhost:${PORT}/ws/interview-room`
    }
  });

  let testExitCode = 0;

  testProcess.stdout.on('data', (data) => {
    const output = data.toString();
    const lines = output.trim().split('\n');
    lines.forEach(line => {
      console.log(`\x1b[32m[TEST FE]\x1b[0m ${line}`);
    });
  });

  testProcess.stderr.on('data', (data) => {
    const output = data.toString();
    const lines = output.trim().split('\n');
    lines.forEach(line => {
      console.error(`\x1b[31m[TEST ERR]\x1b[0m ${line}`);
    });
  });

  await new Promise((resolve) => {
    testProcess.on('close', (code) => {
      testExitCode = code;
      resolve();
    });
  });

  // 3. Cleanup: Tắt Backend server
  console.log('\n[BE] Đang tắt Go Realtime server...');
  server.kill('SIGINT');

  await new Promise((resolve) => {
    server.on('close', () => {
      console.log('[BE] Server đã shutdown thành công.');
      resolve();
    });
  });

  // Xóa file thực thi tạm thời
  try {
    if (fs.existsSync(BINARY_PATH)) {
      fs.unlinkSync(BINARY_PATH);
      console.log('[BE] Đã xóa file thực thi tạm thời.');
    }
  } catch (err) {
    console.warn('[BE] Không thể xóa file thực thi tạm thời:', err.message);
  }

  if (testExitCode === 0) {
    console.log('\n✅ KẾT QUẢ: KIỂM THỬ THÀNH CÔNG! KHÔNG PHÁT HIỆN BUG.');
    process.exit(0);
  } else {
    console.error(`\n❌ KẾT QUẢ: KIỂM THỬ THẤT BẠI! Test exit code: ${testExitCode}`);
    process.exit(1);
  }
}

main().catch(err => {
  console.error('Lỗi chạy kịch bản kiểm thử:', err);
  process.exit(1);
});
