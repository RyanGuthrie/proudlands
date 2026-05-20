<div class="about">
	<div class="hero">
		<h1>Proud Lands</h1>
		<p class="tagline">Making public lands discoverable and usable for everyone.</p>
	</div>

	<section>
		<h2>What is Proud Lands?</h2>
		<p>
			The United States is home to over 640 million acres of public land — national forests,
			wilderness areas, BLM tracts, state parks, and more. These lands belong to every American,
			yet they remain largely invisible to most people. Confusing jurisdiction boundaries, scattered
			data sources, and the lack of a unified tool mean that millions of acres go unexplored and
			underutilized every year.
		</p>
		<p>
			Proud Lands is built to change that. By bringing together authoritative geospatial data from
			sources like the USGS National Map, we aim to create a single, intuitive interface for
			understanding what land is public, where it is, and how to access it.
		</p>
	</section>

	<section>
		<h2>The Mission</h2>
		<p>
			Public lands exist for the benefit of all people — hikers, hunters, anglers, families,
			photographers, and anyone who wants to spend time outside. But discovering and navigating these
			spaces shouldn't require specialized knowledge or hours of research.
		</p>
		<p>
			Proud Lands is built on the belief that access to information about public land is itself a
			form of access to the land. When people can clearly see what's available, where it is, and how
			to get there, they're far more likely to use it, advocate for it, and help protect it.
		</p>
	</section>

	<section>
		<h2>What We're Building</h2>
		<ul>
			<li>An interactive map of U.S. public lands with multiple base layer options</li>
			<li>Trail, road, and access point overlays sourced from open government datasets</li>
			<li>Land status and jurisdiction information at a glance</li>
			<li>Shareable map views so people can send each other exactly where they're going</li>
			<li>Seasonal conditions and access information to help with trip planning</li>
		</ul>
	</section>

	<section class="author">
		<h2>About the Author</h2>
		<p>
			Proud Lands is built by <a href="https://www.linkedin.com/in/mrryanguthrie/" target="_blank" rel="noopener noreferrer">Ryan Guthrie</a>,
			a software engineer and outdoor enthusiast based in Colorado. Ryan built Proud Lands out of
			frustration with how hard it is to find and navigate public land — and a conviction that great
			software can make the outdoors more accessible to everyone.
		</p>
	</section>
</div>

<style>
	.about {
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

	h2 {
		margin: 0;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		color: var(--text-muted);
	}

	p {
		margin: 0;
		line-height: 1.75;
		color: var(--text);
		font-size: 0.95rem;
	}

	ul {
		margin: 0;
		padding-left: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	li {
		line-height: 1.6;
		color: var(--text);
		font-size: 0.95rem;
	}

	.author {
		border-top: 1px solid var(--border);
		padding-top: 2rem;
	}

	a {
		color: var(--accent);
		text-decoration: none;
		font-weight: 500;
	}

	a:hover {
		color: var(--accent-hover);
		text-decoration: underline;
	}
</style>