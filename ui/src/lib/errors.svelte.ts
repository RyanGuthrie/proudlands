export type ApiError = {
	id: number;
	message: string;
	timestamp: Date;
};

let nextId = 1;

export const errors = $state<ApiError[]>([]);

export function reportError(message: string) {
	errors.push({ id: nextId++, message, timestamp: new Date() });
}

type ProblemDetails = { title?: string; detail?: string } | null | undefined;

export function describeApiError(context: string, error?: ProblemDetails): string {
	const reason = error?.detail || error?.title;
	return reason ? `${context}: ${reason}` : context;
}

export function describeNetworkError(context: string, err: unknown): string {
	const reason = err instanceof Error ? err.message : String(err);
	return `${context}: ${reason}`;
}
