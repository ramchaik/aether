import { UserButton } from "@clerk/nextjs";
import { auth, currentUser } from "@clerk/nextjs/server";
import {
  Button,
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
} from "@nextui-org/react";
import Link from "next/link";
import AetherLogo from "./aether-logo";

export default async function Header() {
  const { userId } = auth();
  const user = await currentUser();

  return (
    <Navbar isBordered maxWidth="xl">
      <NavbarBrand>
        <Link href="/" className="font-semibold text-inherit">
          <div className="flex">
            <AetherLogo height={50} width={50} />
            <span className="mt-3">Aether</span>
          </div>
        </Link>
      </NavbarBrand>
      <NavbarContent justify="end">
        {userId && user ? (
          <>
            <NavbarItem>
              <Link href="/dashboard" className="text-inherit">
                Dashboard
              </Link>
            </NavbarItem>
            <NavbarItem>
              <UserButton />
            </NavbarItem>
          </>
        ) : (
          <>
            <NavbarItem>
              <Link href="/sign-in">
                <Button color="primary" variant="flat">
                  Sign In
                </Button>
              </Link>
            </NavbarItem>
            <NavbarItem>
              <Link href="/sign-up">
                <Button color="primary" variant="solid">
                  Sign Up
                </Button>
              </Link>
            </NavbarItem>
          </>
        )}
      </NavbarContent>
    </Navbar>
  );
}
