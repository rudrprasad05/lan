import { Device } from "../lib/models/types";

export type DeviceContextType = {
  devices: Device[];
  error: string | null;
  isLoading: boolean;
};

export type DeviceProviderProps = {
  children: React.ReactNode;
};
