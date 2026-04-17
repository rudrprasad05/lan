export type DeviceType = "laptop" | "desktop" | "mobile";

export type Device = {
  id: string;
  type: DeviceType;
  name: string;
  os: string;
  ip: string;
  ping: string;
};

export type ApiDevice = {
  id: string;
  name: string;
  type: string;
  os: string;
  arch: string;
  ip: string;
  state: string;
  lastSeen: number;
};

export type ApiDevicesResponse = {
  devices: ApiDevice[];
};

export type DeviceCardProps = {
  device: Device;
};
