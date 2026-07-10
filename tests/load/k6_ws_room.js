// k6 WebSocket load/stability test for the realtime interview gateway.
//
// Usage:
//   k6 run tests/load/k6_ws_room.js
//   k6 run -e WS=ws://localhost:8081 -e TOKEN=<jwt> -e VUS=100 tests/load/k6_ws_room.js
//
// Each VU opens a room connection, joins, sends heartbeats, and holds the
// socket open to simulate concurrent interview participants. Thresholds fail
// the run if connections error out or the gateway stops responding.

import ws from 'k6/ws';
import { check } from 'k6';
import { Counter, Rate } from 'k6/metrics';

const WS = __ENV.WS || 'ws://localhost:8081';
const TOKEN = __ENV.TOKEN || 'devtoken';
const VUS = parseInt(__ENV.VUS || '50', 10);
const HOLD_SECONDS = parseInt(__ENV.HOLD || '20', 10);

const wsErrors = new Rate('ws_errors');
const msgsReceived = new Counter('ws_messages_received');

export const options = {
  vus: VUS,
  iterations: VUS,
  thresholds: {
    ws_errors: ['rate<0.05'],
    ws_session_duration: ['p(95)<60000'],
  },
};

export default function () {
  const url = `${WS}/ws/interview-room?token=${TOKEN}`;
  const roomId = `load-room-${__VU % 10}`; // spread across 10 rooms

  const res = ws.connect(url, {}, function (socket) {
    socket.on('open', function () {
      // Join the room.
      socket.send(JSON.stringify({
        event: 'room:join',
        request_id: `join-${__VU}`,
        room_id: roomId,
        payload: { display_name: `VU${__VU}` },
      }));

      // Heartbeat every 15s.
      socket.setInterval(function () {
        socket.send(JSON.stringify({
          event: 'room:heartbeat',
          request_id: `hb-${__VU}-${Date.now()}`,
          room_id: roomId,
          payload: {},
        }));
      }, 15000);

      // Close after the hold window.
      socket.setTimeout(function () {
        socket.close(1000);
      }, HOLD_SECONDS * 1000);
    });

    socket.on('message', function () {
      msgsReceived.add(1);
    });

    socket.on('error', function (e) {
      wsErrors.add(1);
      console.error(`VU${__VU} ws error: ${e.error()}`);
    });
  });

  check(res, { 'ws handshake 101': (r) => r && r.status === 101 }) || wsErrors.add(1);
}
