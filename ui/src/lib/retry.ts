export type CancelRetry = () => void;

/**
 * Repeatedly calls `attempt` until it resolves true, waiting between
 * attempts with exponential backoff (capped at `maxDelayMs`). Retries
 * indefinitely until cancelled or successful.
 */
export function retryWithBackoff(
	attempt: () => Promise<boolean>,
	opts: { initialDelayMs?: number; maxDelayMs?: number } = {}
): CancelRetry {
	const initialDelayMs = opts.initialDelayMs ?? 1000;
	const maxDelayMs = opts.maxDelayMs ?? 5000;
	let cancelled = false;
	let timer: ReturnType<typeof setTimeout> | undefined;

	function scheduleNext(delay: number) {
		timer = setTimeout(async () => {
			if (cancelled) return;
			const success = await attempt();
			if (cancelled || success) return;
			scheduleNext(Math.min(delay * 2, maxDelayMs));
		}, delay);
	}

	scheduleNext(initialDelayMs);

	return () => {
		cancelled = true;
		clearTimeout(timer);
	};
}
