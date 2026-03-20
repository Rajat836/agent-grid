import { redirect } from "next/navigation";

export default function Home() {
  redirect("/system-analysis");
  return null;
}
