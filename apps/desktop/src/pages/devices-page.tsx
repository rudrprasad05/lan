import { Laptop, Monitor, RefreshCcw, Smartphone } from "lucide-react";
import { useState } from "react";
import { SendFileModal } from "../components/file/SendFile";
import { useDevices } from "../context/device-context";
import { DeviceCardProps } from "../lib/models/types";

const API_BASE_URL = "http://127.0.0.1:43821";

export function DevicesPage() {
  const { devices, error, isLoading } = useDevices();

  return (
    <div>
      <PairedDevices />

      <div className="mb-3 flex items-center justify-between">
        <h2 className="font-semibold text-white">• Available Devices</h2>
        <button className="hover:cursor-pointer flex items-center gap-2">
          {isLoading && "Scanning"}
          <RefreshCcw className="w-6 h-6" />
        </button>
      </div>

      {isLoading ? (
        <div className="rounded-lg border border-[#111a2e] bg-[#0b1220] p-4 text-sm text-gray-400">
          Loading devices...
        </div>
      ) : null}

      {error ? (
        <div className="mb-4 rounded-lg border border-red-900 bg-[#0b1220] p-4 text-sm text-red-300">
          {error}
        </div>
      ) : null}

      <div className="grid grid-cols-3 gap-4">
        {devices.map((device) => (
          <DeviceCard key={device.id} device={device} />
        ))}
      </div>

      {!isLoading && devices.length === 0 ? (
        <div className="mt-4 rounded-lg border border-[#111a2e] bg-[#0b1220] p-4 text-sm text-gray-400">
          No LAN devices discovered yet.
        </div>
      ) : null}

      <StatusBar />
    </div>
  );
}

function PairedDevices() {
  return (
    <div className="mb-6 rounded-lg border border-[#111a2e] bg-[#0b1220] p-6">
      <div className="flex flex-col items-center text-gray-400">
        <div className="text-3xl">🔗</div>
        <div className="mt-2 font-semibold text-white">No paired devices</div>
        <div className="text-sm text-gray-400">
          Connect to local devices to start instant sharing.
        </div>
      </div>
    </div>
  );
}

function DeviceCard({ device }: DeviceCardProps) {
  const [open, setOpen] = useState(false);
  const [sendError, setSendError] = useState<string | null>(null);
  const [sendSuccess, setSendSuccess] = useState<string | null>(null);
  const { devices } = useDevices();

  const icon =
    device.type == "laptop" ? (
      <Laptop />
    ) : device.type == "desktop" ? (
      <Monitor />
    ) : (
      <Smartphone />
    );

  const os =
    device.os == "darwin" ? "OSX" : device.os == "windows" ? "WIN" : "Unknown";

  return (
    <div className="flex flex-col gap-3 rounded-lg border border-[#111a2e] bg-[#0b1220] p-4">
      <div className="flex items-center gap-2 text-gray-300">
        {icon}
        <span className="rounded bg-blue-600 px-2 py-0.5 text-xs">{os}</span>
      </div>

      <div>
        <div className="font-medium text-white">{device.name}</div>
        <div className="text-xs text-gray-400">
          {device.ip}
          {device.ping ? ` • ${device.ping}` : ""}
        </div>
      </div>

      <button
        onClick={() => setOpen((prev) => !prev)}
        className="mt-auto rounded bg-blue-600 py-1.5 text-sm text-white hover:bg-blue-500"
      >
        Send
      </button>

      {sendSuccess ? (
        <div className="rounded border border-green-900 bg-green-950/40 px-3 py-2 text-xs text-green-300">
          {sendSuccess}
        </div>
      ) : null}

      {sendError ? (
        <div className="rounded border border-red-900 bg-red-950/40 px-3 py-2 text-xs text-red-300">
          {sendError}
        </div>
      ) : null}

      <SendFileModal
        open={open}
        onClose={() => setOpen(false)}
        devices={devices}
        selectedDeviceId={device.id}
        onSend={async (file, deviceId) => {
          setSendError(null);
          setSendSuccess(null);

          try {
            const formData = new FormData();
            formData.append("file", file);
            formData.append("deviceId", deviceId);

            const response = await fetch(`${API_BASE_URL}/api/files`, {
              method: "POST",
              body: formData,
            });

            if (!response.ok) {
              throw new Error(`File upload failed with ${response.status}.`);
            }

            const payload = (await response.json()) as {
              file?: { name?: string };
            };

            setSendSuccess(
              `Uploaded ${payload.file?.name ?? file.name} for ${device.name}.`,
            );
          } catch (error) {
            const message =
              error instanceof Error ? error.message : "File upload failed.";
            setSendError(message);
            throw error;
          }
        }}
      />
    </div>
  );
}

function StatusBar() {
  return (
    <div className="mt-6 rounded-lg border border-[#111a2e] bg-[#0b1220] p-4 text-white">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-sm">Active Ether Status</div>
          <div className="text-xs text-gray-400">NETWORK SATURATION</div>
        </div>
        <div className="text-sm">0% ACTIVITY</div>
        <div className="text-xs uppercase text-green-400">
          ENCRYPTION ENABLED
        </div>
      </div>

      <div className="mt-3 h-2 rounded bg-gray-800">
        <div className="h-2 w-[3%] rounded bg-blue-600" />
      </div>

      <div className="mt-2 flex justify-between text-xs text-gray-400">
        <div>UPLINK 0.0 Mbps</div>
        <div>DOWNLINK 0.0 Mbps</div>
      </div>

      <div className="mt-3 flex justify-end">
        <div className="flex items-center gap-2 rounded bg-gray-800 px-3 py-1 text-xs">
          <span className="h-2 w-2 rounded-full bg-orange-500" />
          Discovering peers via multicast...
        </div>
      </div>
    </div>
  );
}
