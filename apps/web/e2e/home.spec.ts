import { expect, test } from '@playwright/test';

const sampleRecords = [
  {
    ts: '2026-06-24T10:00:00.000Z',
    session: 'sess-1',
    seq: 1,
    index: 0,
    type: 'request',
    url: '/aiserver.v1.ChatService/StreamChat',
    host: 'api2.cursor.sh',
    direction: 'C2S',
  },
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

test.describe('gRPC Inspector', () => {
  test('loads the main page', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'gRPC Inspector' })).toBeVisible();
  });

  test('displays session count from mocked API', async ({ page }) => {
    await page.route('**/api/records*', async (route) => {
      await route.fulfill({ json: sampleRecords });
    });

    await page.goto('/');
    await expect(page.getByText('2 calls')).toBeVisible({ timeout: 15_000 });
  });
});
