import { ChevronDown, SendHorizonal, Upload, X } from "lucide-react";
import React, { useEffect, useMemo, useRef, useState } from "react";

export type DeviceOption = {
  id: string;
  name: string;
  os: string;
};

type SendFileModalProps = {
  open: boolean;
  onClose: () => void;
  onSend: (file: File, deviceId: string) => Promise<void> | void;
  devices: DeviceOption[];
  selectedDeviceId?: string;
  title?: string;
};

export function SendFileModal({
  open,
  onClose,
  onSend,
  devices,
  selectedDeviceId,
  title = "Send File",
}: SendFileModalProps) {
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [deviceId, setDeviceId] = useState<string>("");
  const [file, setFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isSending, setIsSending] = useState(false);

  useEffect(() => {
    if (!open) return;

    if (selectedDeviceId && devices.some((d) => d.id === selectedDeviceId)) {
      setDeviceId(selectedDeviceId);
      return;
    }

    setDeviceId((prev) => {
      if (prev && devices.some((d) => d.id === prev)) return prev;
      return devices[0]?.id ?? "";
    });
  }, [open, selectedDeviceId, devices]);

  useEffect(() => {
    if (!open) {
      setFile(null);
      setIsDragging(false);
      setIsSending(false);
    }
  }, [open]);

  const selectedDevice = useMemo(
    () => devices.find((d) => d.id === deviceId) ?? null,
    [devices, deviceId],
  );

  const canSend = !!file && !!deviceId && !isSending;

  const handleBrowseClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const nextFile = event.target.files?.[0] ?? null;
    if (nextFile) {
      setFile(nextFile);
    }
  };

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    event.stopPropagation();
    setIsDragging(false);

    const nextFile = event.dataTransfer.files?.[0] ?? null;
    if (nextFile) {
      setFile(nextFile);
    }
  };

  const handleSend = async () => {
    if (!file || !deviceId) return;

    try {
      setIsSending(true);
      await onSend(file, deviceId);
      onClose();
    } finally {
      setIsSending(false);
    }
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-[#031230]/75 backdrop-blur-md">
      <div className="w-full max-w-140 overflow-hidden rounded-2xl border border-white/10 bg-[#202944] shadow-2xl">
        <div className="flex items-center justify-between px-7 pt-7 pb-5">
          <div className="flex items-center gap-3">
            <SendHorizonal className="h-6 w-6 text-sky-300" />
            <h2 className="text-4 font-semibold text-white">{title}</h2>
          </div>

          <button
            type="button"
            onClick={onClose}
            className="rounded-md p-1 text-slate-400 transition hover:bg-white/5 hover:text-white"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="px-7 pb-6">
          <div className="space-y-8">
            <div className="space-y-3">
              <label className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">
                Recipient
              </label>

              <div className="relative">
                <select
                  value={deviceId}
                  onChange={(e) => setDeviceId(e.target.value)}
                  className="h-14 w-full appearance-none rounded-xl border border-white/5 bg-[#313A58] px-4 pr-12 text-[22px] text-slate-100 outline-none transition focus:border-sky-400/40"
                >
                  {devices.length === 0 ? (
                    <option value="">No devices found</option>
                  ) : null}

                  {devices.map((device) => (
                    <option key={device.id} value={device.id}>
                      {device.name} ({device.os})
                    </option>
                  ))}
                </select>

                <ChevronDown className="pointer-events-none absolute top-1/2 right-4 h-5 w-5 -translate-y-1/2 text-slate-400" />
              </div>
            </div>

            <div className="space-y-3">
              <label className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">
                File to transfer
              </label>

              <div
                onDragEnter={(e) => {
                  e.preventDefault();
                  setIsDragging(true);
                }}
                onDragOver={(e) => {
                  e.preventDefault();
                  setIsDragging(true);
                }}
                onDragLeave={(e) => {
                  e.preventDefault();
                  setIsDragging(false);
                }}
                onDrop={handleDrop}
                className={`rounded-xl border-2 border-dashed px-6 py-10 text-center transition ${
                  isDragging
                    ? "border-sky-400 bg-sky-500/10"
                    : "border-[#16243e] bg-[#04122B]"
                }`}
              >
                <input
                  ref={fileInputRef}
                  type="file"
                  className="hidden"
                  onChange={handleFileChange}
                />

                {!file ? (
                  <div className="flex flex-col items-center">
                    <div className="mb-6 flex h-20 w-20 items-center justify-center rounded-2xl bg-[#102344]">
                      <Upload className="h-9 w-9 text-sky-200" />
                    </div>

                    <p className="text-2 font-medium text-white">
                      Drag and drop files here
                    </p>
                    <p className="mt-1 text-sm text-slate-400">
                      Maximum file size: 10GB
                    </p>

                    <button
                      type="button"
                      onClick={handleBrowseClick}
                      className="mt-6 rounded-lg bg-[#2E3B63] px-6 py-3 text-sm font-semibold text-slate-200 transition hover:bg-[#394972]"
                    >
                      Browse Files
                    </button>
                  </div>
                ) : (
                  <div className="flex flex-col items-center">
                    <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-[#102344]">
                      <Upload className="h-7 w-7 text-sky-200" />
                    </div>

                    <p className="max-w-full truncate text-base font-medium text-white">
                      {file.name}
                    </p>
                    <p className="mt-1 text-sm text-slate-400">
                      {(file.size / 1024 / 1024).toFixed(2)} MB
                    </p>

                    <button
                      type="button"
                      onClick={handleBrowseClick}
                      className="mt-5 rounded-lg bg-[#2E3B63] px-5 py-2.5 text-sm font-semibold text-slate-200 transition hover:bg-[#394972]"
                    >
                      Change File
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="mt-8 flex items-center justify-end gap-4">
            <button
              type="button"
              onClick={onClose}
              className="text-base font-medium text-slate-300 transition hover:text-white"
            >
              Cancel
            </button>

            <button
              type="button"
              disabled={!canSend}
              onClick={handleSend}
              className="rounded-xl bg-sky-400 px-8 py-3 text-base font-semibold text-[#06223A] shadow-[0_0_24px_rgba(56,189,248,0.35)] transition hover:bg-sky-300 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {isSending ? "Sending..." : "Send"}
            </button>
          </div>
        </div>

        <div className="flex items-center justify-between border-t border-white/5 bg-[#07142C] px-7 py-3 text-xs uppercase tracking-[0.12em] text-slate-400">
          <span>Direct LAN encrypted</span>
          <span>
            {selectedDevice
              ? `Ready for ${selectedDevice.name}`
              : "No recipient"}
          </span>
        </div>
      </div>
    </div>
  );
}
