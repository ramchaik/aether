"use client";

import { NextUIProvider } from "@nextui-org/react";
import { ClerkProvider } from "@clerk/nextjs";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <NextUIProvider>
      <ClerkProvider>{children}</ClerkProvider>
    </NextUIProvider>
  );
}
