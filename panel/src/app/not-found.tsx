import Image from "next/image";
import Link from "next/link";
import { SidebarLayout } from "@/components/sidebar-layout";
import { Button } from "@/components/ui/button";

export default function NotFound() {
  return (
    <SidebarLayout>
      <div className="flex flex-1 flex-col items-center justify-center gap-6 px-6 py-[14px] text-center">
        <Image
          src="/error-404.jpg"
          alt="404 Not Found"
          width={600}
          height={600}
          className="object-contain"
          priority
        />
        <Link href="/">
          <Button className="mt-4 cursor-pointer">Return Home</Button>
        </Link>
      </div>
    </SidebarLayout>
  );
}
