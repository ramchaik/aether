import { UserButton } from "@clerk/nextjs";
import { auth, currentUser } from "@clerk/nextjs/server";
import Link from "next/link";
import React from "react";

export default async function Header() {
  const { userId } = auth();
  const user = await currentUser();
  return (
    <div className="bg-gray-600 text-neutral-100">
      <div className="container mx-auto flex items-center justify-between py-4">
        <Link href="/">Home</Link>
        <div>
          <div className="flex gap-4 items-center">
            {userId && user ? (
              <>
                <Link href="/dashboard">Dashboard</Link>
                <UserButton afterSignOutUrl="/" />
              </>
            ) : (
              <>
                <Link href="/sign-in">Sign in</Link>
                <Link href="/sign-up">Sign up</Link>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
