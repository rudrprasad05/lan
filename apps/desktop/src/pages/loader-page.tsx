type LoaderPageProps = {
  error: string | null;
  isLoading: boolean;
  onRetry: () => void;
};

export function LoaderPage({ error, isLoading, onRetry }: LoaderPageProps) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-black px-6 text-white">
      <div className="w-full max-w-md rounded-lg border border-[#111a2e] bg-[#0b1220] p-6">
        <div className="text-lg font-semibold">Starting LAN desktop</div>
        <div className="mt-2 text-sm text-gray-400">
          Checking backend health before loading the app.
        </div>

        <div className="mt-6 rounded border border-[#111a2e] bg-[#050914] px-4 py-3 text-sm">
          {isLoading ? "Probing http://127.0.0.1:43821/api/health..." : error}
        </div>

        {error ? (
          <button
            className="mt-4 rounded bg-blue-600 px-4 py-2 text-sm text-white hover:bg-blue-500"
            onClick={onRetry}
            type="button"
          >
            Retry
          </button>
        ) : null}
      </div>
    </div>
  );
}
