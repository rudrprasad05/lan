import { useEffect, useState } from "react";
import { AppShell } from "./components/app-shell";
import { DeviceProvider } from "./context/device-context";
import { checkBackendHealth } from "./lib/api";
import { DashboardPage } from "./pages/dashboard-page";
import { DevicesPage } from "./pages/devices-page";
import { LoaderPage } from "./pages/loader-page";
import { SharePage } from "./pages/share-page";

export type AppPage = "dashboard" | "devices" | "share";

export default function App() {
  const [isBootstrapping, setIsBootstrapping] = useState(true);
  const [healthError, setHealthError] = useState<string | null>(null);
  const [page, setPage] = useState<AppPage>("devices");

  useEffect(() => {
    let cancelled = false;

    async function bootstrap() {
      setIsBootstrapping(true);
      setHealthError(null);

      try {
        const res = await checkBackendHealth();
        console.log(res);
        if (cancelled) return;
        setPage("devices");
      } catch (error) {
        if (cancelled) return;
        setHealthError(
          error instanceof Error
            ? error.message
            : "Backend health check failed.",
        );
        setIsBootstrapping(false);
        return;
      }

      if (!cancelled) {
        setIsBootstrapping(false);
      }
    }

    bootstrap();

    return () => {
      cancelled = true;
    };
  }, []);

  if (isBootstrapping || healthError) {
    return (
      <LoaderPage
        error={healthError}
        isLoading={isBootstrapping}
        onRetry={() => {
          setIsBootstrapping(true);
          setHealthError(null);
          checkBackendHealth()
            .then(() => {
              setPage("devices");
              setIsBootstrapping(false);
            })
            .catch((error: unknown) => {
              setHealthError(
                error instanceof Error
                  ? error.message
                  : "Backend health check failed.",
              );
              setIsBootstrapping(false);
            });
        }}
      />
    );
  }

  return (
    <DeviceProvider>
      <AppShell currentPage={page} onNavigate={setPage}>
        {page === "dashboard" ? <DashboardPage /> : null}
        {page === "devices" ? <DevicesPage /> : null}
        {page === "share" ? <SharePage /> : null}
      </AppShell>
    </DeviceProvider>
  );
}
