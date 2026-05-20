<script lang="ts">
	import treadLightlyBadge from '$lib/assets/tread-lightly_badge-primary.svg';
</script>

<div class="home">
	<div class="hero">
		<h1>Proud Lands</h1>
		<p class="tagline">Your public lands, made discoverable.</p>
	</div>

	<section>
		<p>
			Over 640 million acres of land in the United States are publicly owned — and most people have
			no idea what's out there or how to reach it. Proud Lands is a mapping tool built to change
			that, putting authoritative data about national forests, wilderness areas, BLM tracts, and more
			into a single, easy-to-use interface.
		</p>
		<p>
			Whether you're planning a weekend trip, scouting a new trail, or just curious what's beyond
			the road you drive every day, Proud Lands helps you see what belongs to everyone.
		</p>
	</section>

	<div class="cta-row">
		<a href="/map" class="cta-primary">Open the Map</a>
		<a href="/about" class="cta-secondary">Learn More</a>
	</div>

	<section class="tread-lightly">
		<a href="https://treadlightly.org" target="_blank" rel="noopener noreferrer" class="tread-badge">
			<img src={treadLightlyBadge} alt="Tread Lightly!" />
		</a>
		<div>
			<h2>Tread Lightly</h2>
			<p>
				Being proud of our public lands means more than enjoying them — it means protecting,
				respecting, and fostering them for the generations that follow. Leave no trace, stay on
				designated routes, and treat every acre as the shared treasure it is. Proud Lands is proud
				to support the principles of <a href="https://treadlightly.org" target="_blank" rel="noopener noreferrer">Tread Lightly!</a>,
				a nonprofit dedicated to promoting responsible outdoor recreation on public and private land.
			</p>
		</div>
	</section>
</div>

<style>
	.home {
		max-width: 680px;
		margin: 0 auto;
		padding: 3rem 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 2.5rem;
	}

	.hero {
		border-bottom: 1px solid var(--border);
		padding-bottom: 2rem;
	}

	.hero h1 {
		margin: 0 0 0.5rem;
		font-size: 2.2rem;
		font-weight: 800;
		color: var(--text);
		letter-spacing: -0.02em;
	}

	.tagline {
		margin: 0;
		font-size: 1.1rem;
		color: var(--accent);
		font-weight: 500;
	}

	section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	p {
		margin: 0;
		line-height: 1.75;
		color: var(--text);
		font-size: 0.95rem;
	}

	.cta-row {
		display: flex;
		gap: 0.75rem;
		align-items: center;
	}

	.cta-primary {
		background: var(--accent);
		color: #0f1117;
		font-weight: 600;
		font-size: 0.9rem;
		padding: 0.55rem 1.25rem;
		border-radius: 6px;
		text-decoration: none;
		transition: background 0.15s;
	}

	.cta-primary:hover {
		background: var(--accent-hover);
	}

	.cta-secondary {
		color: var(--text-muted);
		font-size: 0.9rem;
		text-decoration: none;
		transition: color 0.15s;
	}

	.cta-secondary:hover {
		color: var(--text);
	}

	.tread-lightly {
		display: flex;
		gap: 1.5rem;
		align-items: flex-start;
		border-top: 1px solid var(--border);
		padding-top: 2rem;
	}

	.tread-badge {
		flex-shrink: 0;
	}

	.tread-badge img {
		width: 90px;
		height: auto;
		display: block;
		opacity: 0.9;
		transition: opacity 0.15s;
	}

	.tread-badge:hover img {
		opacity: 1;
	}

	.tread-lightly h2 {
		margin: 0 0 0.6rem;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		color: var(--text-muted);
	}

	.tread-lightly a {
		color: var(--accent);
		text-decoration: none;
		font-weight: 500;
	}

	.tread-lightly a:hover {
		color: var(--accent-hover);
		text-decoration: underline;
	}
</style>