import { ApiDevicesResponse, Device } from "./models/types";

const API_BASE_URL = "http://127.0.0.1:43821";

export async function checkBackendHealth() {
  const response = await fetch(`${API_BASE_URL}/api/health`);
  if (!response.ok) {
    throw new Error(`Backend health check failed with ${response.status}.`);
  }
  const data = await response.json();
  return data;
}

export async function getDevices() {
  const response = await fetch(`${API_BASE_URL}/api/devices`);
  if (!response.ok) {
    throw new Error(`Device fetch failed with ${response.status}.`);
  }

  const payload = (await response.json()) as ApiDevicesResponse;
  return payload.devices.map(mapApiDevice);
}

function mapApiDevice(device: ApiDevicesResponse["devices"][number]): Device {
  return {
    id: device.id,
    ip: device.ip,
    name: device.name,
    os: device.os,
    ping: formatLastSeen(device.lastSeen),
    type: normalizeDeviceType(device.type, device.os),
  };
}

function normalizeDeviceType(type: string, os: string): Device["type"] {
  const normalizedType = type.toLowerCase();
  const normalizedOs = os.toLowerCase();

  if (normalizedType === "laptop") {
    return "laptop";
  }

  if (normalizedType === "desktop") {
    return "desktop";
  }

  if (
    normalizedType === "mobile" ||
    normalizedType === "phone" ||
    normalizedOs.includes("android") ||
    normalizedOs.includes("ios")
  ) {
    return "mobile";
  }

  return "desktop";
}

function formatLastSeen(lastSeen: number) {
  const ageSeconds = Math.max(0, Math.floor(Date.now() / 1000) - lastSeen);
  return ageSeconds <= 1 ? "seen just now" : `seen ${ageSeconds}s ago`;
}
