import {
  Grid,
  Laptop,
  RefreshCcw,
  Settings,
  Share2,
} from "lucide-react";
import { ReactNode } from "react";
import { AppPage } from "../App";

type AppShellProps = {
  children: ReactNode;
  currentPage: AppPage;
  onNavigate: (page: AppPage) => void;
};

type NavItemProps = {
  active: boolean;
  icon: ReactNode;
  label: string;
  onClick: () => void;
};

export function AppShell({
  children,
  currentPage,
  onNavigate,
}: AppShellProps) {
  return (
    <div className="flex h-screen bg-black">
      <div className="flex w-64 flex-col justify-between border-r border-[#111a2e] bg-[#0b1220] text-white">
        <div>
          <div className="p-4">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-600">
                ⚡
              </div>
              <div>
                <div className="font-semibold">Kinetic Ether</div>
                <div className="text-xs text-blue-300">NETWORK ACTIVE</div>
              </div>
            </div>
          </div>

          <nav className="mt-6 space-y-1 px-3">
            <NavItem
              active={currentPage === "dashboard"}
              icon={<Grid size={18} />}
              label="Dashboard"
              onClick={() => onNavigate("dashboard")}
            />
            <NavItem
              active={currentPage === "devices"}
              icon={<Laptop size={18} />}
              label="Devices"
              onClick={() => onNavigate("devices")}
            />
            <NavItem
              active={currentPage === "share"}
              icon={<Share2 size={18} />}
              label="Share"
              onClick={() => onNavigate("share")}
            />
          </nav>
        </div>

        <div className="space-y-3 p-3">
          <button className="w-full rounded bg-cyan-500 py-2 font-medium text-black hover:bg-cyan-400">
            → Send File
          </button>
          <div className="flex items-center gap-2 px-2 text-sm text-gray-400">
            <Settings size={16} /> Settings
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-linear-to-br from-[#050914] to-[#0a0f1f] p-6">
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-2xl font-semibold text-white">
            {currentPage === "devices"
              ? "Network Devices"
              : currentPage === "dashboard"
                ? "Dashboard"
                : "Share"}
          </h1>
          <div className="flex items-center gap-3">
            <RefreshCcw className="text-gray-400" size={18} />
            <div className="h-8 w-8 rounded-full bg-gray-600" />
          </div>
        </div>

        {children}
      </div>
    </div>
  );
}

function NavItem({ active, icon, label, onClick }: NavItemProps) {
  return (
    <button
      className={`flex w-full cursor-pointer items-center gap-3 rounded px-3 py-2 text-left ${
        active ? "bg-blue-600 text-white" : "text-gray-300 hover:bg-[#111a2e]"
      }`}
      onClick={onClick}
      type="button"
    >
      {icon}
      <span>{label}</span>
    </button>
  );
}
