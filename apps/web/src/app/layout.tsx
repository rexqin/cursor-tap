import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "gRPC Inspector",
  description: "Real-time gRPC traffic inspector for MITM proxy",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased">{children}</body>
    </html>
  );
}
