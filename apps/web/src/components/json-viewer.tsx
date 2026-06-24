'use client';

import { useMemo, useState } from 'react';
import { Copy, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface JsonViewerProps {
  data: string;
  maxHeight?: string;
}

const MAX_DISPLAY_CHARS = 100_000;
const MAX_HIGHLIGHT_CHARS = 20_000;

export function JsonViewer({ data, maxHeight = '400px' }: JsonViewerProps) {
  const [copied, setCopied] = useState(false);

  const { formatted, truncated } = useMemo(() => {
    let text = data;
    let truncated = false;

    if (text.length > MAX_DISPLAY_CHARS) {
      text = text.slice(0, MAX_DISPLAY_CHARS);
      truncated = true;
    }

    try {
      const parsed = JSON.parse(text);
      return { formatted: JSON.stringify(parsed, null, 2), truncated };
    } catch {
      return { formatted: text, truncated };
    }
  }, [data]);

  const highlighted = useMemo(() => {
    if (formatted.length > MAX_HIGHLIGHT_CHARS) {
      return null;
    }
    return formatted
      .replace(/"([^"]+)":/g, '<span class="text-purple-600 dark:text-purple-400">"$1"</span>:')
      .replace(/: "([^"]*)"/g, ': <span class="text-green-600 dark:text-green-400">"$1"</span>')
      .replace(/: (\d+)/g, ': <span class="text-blue-600 dark:text-blue-400">$1</span>')
      .replace(/: (true|false)/g, ': <span class="text-orange-600 dark:text-orange-400">$1</span>')
      .replace(/: (null)/g, ': <span class="text-gray-500">$1</span>');
  }, [formatted]);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(data);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  return (
    <div className="relative group w-full min-w-0">
      <Button
        variant="ghost"
        size="sm"
        className="absolute top-2 right-2 h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity z-10"
        onClick={handleCopy}
        title="Copy to clipboard"
      >
        {copied ? (
          <Check className="w-4 h-4 text-green-500" />
        ) : (
          <Copy className="w-4 h-4" />
        )}
      </Button>
      {truncated && (
        <div className="mb-1 text-xs text-amber-600 dark:text-amber-400">
          内容过大，仅显示前 {MAX_DISPLAY_CHARS.toLocaleString()} 字符（复制仍可获取完整内容）
        </div>
      )}
      <div className="w-full overflow-x-auto">
        {highlighted ? (
          <pre
            className="p-3 bg-muted rounded text-xs font-mono whitespace-pre"
            style={{ maxHeight, minWidth: 'min-content' }}
            dangerouslySetInnerHTML={{ __html: highlighted }}
          />
        ) : (
          <pre
            className="p-3 bg-muted rounded text-xs font-mono whitespace-pre overflow-auto"
            style={{ maxHeight, minWidth: 'min-content' }}
          >
            {formatted}
          </pre>
        )}
      </div>
    </div>
  );
}
