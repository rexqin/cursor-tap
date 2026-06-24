import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { WSClient } from './ws-client';
import { sampleGrpcRecord } from '@/test/fixtures/records';

class MockWebSocket {
  static OPEN = 1;
  static CLOSED = 3;
  static instances: MockWebSocket[] = [];

  readyState = MockWebSocket.CLOSED;
  url: string;
  onopen: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onclose: (() => void) | null = null;
  onerror: ((error: unknown) => void) | null = null;

  constructor(url: string) {
    this.url = url;
    MockWebSocket.instances.push(this);
    queueMicrotask(() => {
      this.readyState = MockWebSocket.OPEN;
      this.onopen?.();
    });
  }

  close() {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.();
  }

  simulateMessage(data: string) {
    this.onmessage?.({ data });
  }

  static latest(): MockWebSocket {
    return MockWebSocket.instances[MockWebSocket.instances.length - 1];
  }

  static reset() {
    MockWebSocket.instances = [];
  }
}

describe('WSClient', () => {
  beforeEach(() => {
    vi.stubGlobal('WebSocket', MockWebSocket);
    MockWebSocket.reset();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it('calls onStatus(true) when connected', async () => {
    const onStatus = vi.fn();
    const client = new WSClient('ws://localhost:9090/ws/records', vi.fn(), onStatus);
    client.connect();
    await vi.waitFor(() => expect(onStatus).toHaveBeenCalledWith(true));
  });

  it('parses messages and calls onRecord', async () => {
    const onRecord = vi.fn();
    const client = new WSClient('ws://localhost:9090/ws/records', onRecord, vi.fn());
    client.connect();
    await vi.waitFor(() => expect(MockWebSocket.instances).toHaveLength(1));

    MockWebSocket.latest().simulateMessage(JSON.stringify(sampleGrpcRecord));
    expect(onRecord).toHaveBeenCalledWith(sampleGrpcRecord);
  });

  it('deduplicates are handled by consumer; invalid JSON is ignored', async () => {
    const onRecord = vi.fn();
    const client = new WSClient('ws://localhost:9090/ws/records', onRecord, vi.fn());
    client.connect();
    await vi.waitFor(() => expect(MockWebSocket.instances).toHaveLength(1));

    MockWebSocket.latest().simulateMessage('not-json');
    expect(onRecord).not.toHaveBeenCalled();
  });

  it('schedules reconnect after disconnect', async () => {
    const onStatus = vi.fn();
    const client = new WSClient('ws://localhost:9090/ws/records', vi.fn(), onStatus);
    client.connect();
    await vi.waitFor(() => expect(onStatus).toHaveBeenCalledWith(true));

    MockWebSocket.latest().close();
    expect(onStatus).toHaveBeenCalledWith(false);

    vi.advanceTimersByTime(1000);
    expect(MockWebSocket.instances).toHaveLength(2);
  });

  it('calls onReconnect on second successful connect', async () => {
    const onReconnect = vi.fn();
    const client = new WSClient('ws://localhost:9090/ws/records', vi.fn(), vi.fn(), onReconnect);
    client.connect();
    await vi.waitFor(() => expect(MockWebSocket.instances).toHaveLength(1));

    MockWebSocket.latest().close();
    vi.advanceTimersByTime(1000);
    await vi.waitFor(() => expect(MockWebSocket.instances).toHaveLength(2));

    expect(onReconnect).toHaveBeenCalledTimes(1);
  });

  it('disconnect prevents further reconnects', async () => {
    const client = new WSClient('ws://localhost:9090/ws/records', vi.fn(), vi.fn());
    client.connect();
    await vi.waitFor(() => expect(MockWebSocket.instances).toHaveLength(1));

    client.disconnect();
    MockWebSocket.latest().close();
    vi.advanceTimersByTime(5000);
    expect(MockWebSocket.instances).toHaveLength(1);
  });
});
