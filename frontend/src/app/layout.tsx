import './globals.css'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'ffprobe online',
  description: 'run ffprobe online',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
