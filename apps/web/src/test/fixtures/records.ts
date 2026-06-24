import { Record } from '@/lib/types';

export const sampleRecord: Record = {
  ts: '2026-06-24T10:00:00.000Z',
  session: 'sess-1',
  seq: 1,
  index: 0,
  type: 'request',
  method: 'POST',
  url: '/aiserver.v1.ChatService/StreamChat',
  host: 'api2.cursor.sh',
  direction: 'C2S',
  size: 128,
};

export const sampleGrpcRecord: Record = {
  ts: '2026-06-24T10:00:01.000Z',
  session: 'sess-1',
  seq: 1,
  index: 1,
  type: 'grpc',
  host: 'api2.cursor.sh',
  direction: 'C2S',
  grpc_service: 'aiserver.v1.ChatService',
  grpc_method: 'StreamChat',
  grpc_data: '{"message":"hello"}',
  size: 64,
};

export const sampleRecords: Record[] = [
  sampleRecord,
  sampleGrpcRecord,
  {
    ts: '2026-06-24T10:01:00.000Z',
    session: 'sess-2',
    seq: 2,
    index: 0,
    type: 'request',
    url: '/aiserver.v1.AgentService/Run',
    host: 'api2.cursor.sh',
    direction: 'C2S',
  },
];
