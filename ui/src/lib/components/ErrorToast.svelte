<script lang="ts">
	import { fade } from 'svelte/transition';
	import { errors } from '$lib/errors.svelte';

	const FADE_AFTER_MS = 10_000;

	let expanded = $state(false);
	let visible = $state(false);
	let fadeTimer: ReturnType<typeof setTimeout> | undefined;
	let seenCount = 0;

	function scheduleFade() {
		clearTimeout(fadeTimer);
		fadeTimer = setTimeout(() => {
			visible = false;
		}, FADE_AFTER_MS);
	}

	$effect(() => {
		if (errors.length > seenCount) {
			seenCount = errors.length;
			visible = true;
			if (!expanded) scheduleFade();
		}
	});

	function toggleExpanded() {
		expanded = !expanded;
		if (expanded) {
			clearTimeout(fadeTimer);
		} else {
			scheduleFade();
		}
	}

	function formatTime(d: Date) {
		return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
	}

	let latestError = $derived(errors.at(-1));
</script>

{#if visible && errors.length > 0}
	<div class="error-toast" class:expanded transition:fade={{ duration: 250 }}>
		{#if !expanded}
			<div class="toast-row">
				<span class="toast-dot"></span>
				<span class="toast-message">{latestError?.message}</span>
				<button class="toast-btn" onclick={toggleExpanded}>
					{errors.length > 1 ? `+${errors.length - 1} more` : 'Details'}
				</button>
			</div>
		{:else}
			<div class="toast-panel">
				<div class="toast-panel-header">
					<span>{errors.length} error{errors.length === 1 ? '' : 's'}</span>
					<button class="toast-btn" onclick={toggleExpanded}>Collapse</button>
				</div>
				<ul class="toast-list">
					{#each [...errors].reverse() as err (err.id)}
						<li>
							<span class="toast-time">{formatTime(err.timestamp)}</span>
							<span class="toast-text">{err.message}</span>
						</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>
{/if}

<style>
	.error-toast {
		position: fixed;
		bottom: 1rem;
		left: 50%;
		transform: translateX(-50%);
		z-index: 100;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 8px;
		box-shadow: 0 4px 16px rgba(0, 0, 0, 0.45);
		max-width: min(90vw, 480px);
	}

	.toast-row {
		display: flex;
		align-items: center;
		gap: 0.6rem;
		padding: 0.5rem 0.6rem 0.5rem 0.9rem;
	}

	.toast-dot {
		width: 7px;
		height: 7px;
		border-radius: 50%;
		background: #e57373;
		flex-shrink: 0;
	}

	.toast-message {
		font-size: 0.82rem;
		color: var(--text);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		min-width: 0;
	}

	.toast-btn {
		flex-shrink: 0;
		background: transparent;
		border: 1px solid var(--border);
		border-radius: 4px;
		color: var(--text-muted);
		cursor: pointer;
		font-family: inherit;
		font-size: 0.72rem;
		padding: 0.25rem 0.5rem;
		white-space: nowrap;
		transition: color 0.15s, border-color 0.15s;
	}

	.toast-btn:hover {
		color: var(--text);
		border-color: var(--text-muted);
	}

	.toast-panel {
		width: min(90vw, 420px);
		max-height: 320px;
		display: flex;
		flex-direction: column;
	}

	.toast-panel-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.6rem 0.75rem;
		border-bottom: 1px solid var(--border);
		font-size: 0.78rem;
		font-weight: 600;
		color: var(--text);
	}

	.toast-list {
		list-style: none;
		margin: 0;
		padding: 0.3rem 0;
		overflow-y: auto;
	}

	.toast-list li {
		display: flex;
		gap: 0.6rem;
		padding: 0.4rem 0.75rem;
		border-bottom: 1px solid var(--border);
	}

	.toast-list li:last-child {
		border-bottom: none;
	}

	.toast-time {
		flex-shrink: 0;
		font-size: 0.7rem;
		color: var(--text-muted);
		font-variant-numeric: tabular-nums;
		padding-top: 0.1rem;
	}

	.toast-text {
		font-size: 0.78rem;
		color: var(--text);
		line-height: 1.4;
	}
</style>
