import { createContext, useContext, useEffect, useState } from "react";
import { getDevices } from "../lib/api";
import { Device } from "../lib/models/types";
import { DeviceContextType, DeviceProviderProps } from "./types";

const DeviceContext = createContext<DeviceContextType | undefined>(undefined);

export function DeviceProvider({ children }: DeviceProviderProps) {
  const [devices, setDevices] = useState<Device[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let active = true;

    async function loadDevices() {
      try {
        const nextDevices = await getDevices();
        if (!active) return;
        setDevices(nextDevices);
        setError(null);
      } catch (loadError) {
        if (!active) return;
        setError(
          loadError instanceof Error
            ? loadError.message
            : "Failed to load devices.",
        );
      } finally {
        if (active) {
          setIsLoading(false);
        }
      }
    }

    loadDevices();
    const interval = window.setInterval(loadDevices, 5000);

    return () => {
      active = false;
      window.clearInterval(interval);
    };
  }, []);

  return (
    <DeviceContext.Provider value={{ devices, error, isLoading }}>
      {children}
    </DeviceContext.Provider>
  );
}

export function useDevices() {
  const ctx = useContext(DeviceContext);
  if (!ctx) throw new Error("useDevices must be used inside DeviceProvider");
  return ctx;
}
