<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import maplibregl, { type StyleSpecification } from 'maplibre-gl';
	import 'maplibre-gl/dist/maplibre-gl.css';

	let mapContainer: HTMLDivElement;
	let map: maplibregl.Map;

	const layers = [
		{
			id: 'usgs-topo',
			label: 'USGS Topo',
			tiles: ['https://basemap.nationalmap.gov/arcgis/rest/services/USGSTopo/MapServer/tile/{z}/{y}/{x}'],
			attribution: 'USGS National Map',
		},
		{
			id: 'usgs-imagery',
			label: 'USGS Imagery',
			tiles: ['https://basemap.nationalmap.gov/arcgis/rest/services/USGSImageryOnly/MapServer/tile/{z}/{y}/{x}'],
			attribution: 'USGS National Map',
		},
		{
			id: 'osm',
			label: 'OpenStreetMap',
			tiles: ['https://tile.openstreetmap.org/{z}/{x}/{y}.png'],
			attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
		},
		{
			id: 'stamen-terrain',
			label: 'Stamen Terrain',
			tiles: ['https://tiles.stadiamaps.com/tiles/stamen_terrain/{z}/{x}/{y}.png'],
			attribution: '© <a href="https://stadiamaps.com/">Stadia Maps</a> © <a href="https://stamen.com/">Stamen Design</a> © <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
		},
	] as const;

	type LayerId = (typeof layers)[number]['id'];

	let activeLayerId = $state<LayerId>('usgs-topo');
	let sidebarOpen = $state(true);

	let centerLat = $state(44.0);
	let centerLng = $state(-105.5);
	let zoom = $state(6);

	function fmt(val: number, decimals: number) {
		return val.toFixed(decimals);
	}

	function makeStyle(layerId: LayerId): StyleSpecification {
		const layer = layers.find((l) => l.id === layerId)!;
		return {
			version: 8,
			sources: {
				tiles: {
					type: 'raster',
					tiles: [...layer.tiles],
					tileSize: 256,
					attribution: layer.attribution,
				},
			},
			layers: [{ id: 'tiles', type: 'raster', source: 'tiles' }],
		};
	}

	function switchLayer(id: LayerId) {
		activeLayerId = id;
		map.setStyle(makeStyle(id));
	}

	onMount(() => {
		map = new maplibregl.Map({
			container: mapContainer,
			style: makeStyle(activeLayerId),
			center: [-105.5, 44.0],
			zoom: 6,
		});
		map.addControl(new maplibregl.NavigationControl(), 'bottom-right');
		map.addControl(new maplibregl.ScaleControl(), 'bottom-left');

		function syncPosition() {
			const c = map.getCenter();
			centerLat = c.lat;
			centerLng = c.lng;
			zoom = map.getZoom();
		}

		map.on('move', syncPosition);
	});

	onDestroy(() => {
		map?.remove();
	});
</script>

<div class="map-page">
	<!-- Header -->
	<header class="map-header">
		<div class="header-left">
			<h1>Proud Lands</h1>
			<span class="header-subtitle">Explore public lands, trails, and wilderness areas</span>
		</div>
	</header>

	<!-- Body: sidebar + map -->
	<div class="map-body">
		<!-- Layer switcher: always visible, floating on map -->
		<div class="layer-switcher">
			<span class="layer-switcher-label">Base Layer</span>
			<div class="layer-list">
				{#each layers as layer}
					<button
						class="layer-btn"
						class:active={activeLayerId === layer.id}
						onclick={() => switchLayer(layer.id)}
					>
						{layer.label}
					</button>
				{/each}
			</div>
		</div>

		<!-- Sidebar -->
		<aside class="sidebar" class:collapsed={!sidebarOpen}>
			<div class="sidebar-content">

				<section class="sidebar-section">
					<h2>Location</h2>
					<dl class="data-list">
						<dt>Latitude</dt>
						<dd>{fmt(Math.abs(centerLat), 4)}° {centerLat >= 0 ? 'N' : 'S'}</dd>
						<dt>Longitude</dt>
						<dd>{fmt(Math.abs(centerLng), 4)}° {centerLng >= 0 ? 'E' : 'W'}</dd>
						<dt>Zoom</dt>
						<dd>{fmt(zoom, 1)}</dd>
					</dl>
				</section>

				<section class="sidebar-section">
					<h2>Area Info</h2>
					<dl class="data-list">
						<dt>Region</dt>
						<dd>Greater Yellowstone</dd>
						<dt>Land Status</dt>
						<dd>National Forest</dd>
						<dt>Nearest Town</dt>
						<dd>Cody, WY</dd>
					</dl>
				</section>

				<section class="sidebar-section">
					<h2>Conditions</h2>
					<dl class="data-list">
						<dt>Season</dt>
						<dd>Spring</dd>
						<dt>Snow Level</dt>
						<dd>8,500 ft</dd>
						<dt>Fire Danger</dt>
						<dd>Low</dd>
						<dt>Road Access</dt>
						<dd>Open</dd>
					</dl>
				</section>

			</div>

			<!-- Toggle tab -->
			<button class="sidebar-toggle" onclick={() => (sidebarOpen = !sidebarOpen)} aria-label="Toggle sidebar">
				<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
					{#if sidebarOpen}
						<path d="M11 8L5 3v10l6-5z"/>
					{:else}
						<path d="M5 8l6-5v10L5 8z"/>
					{/if}
				</svg>
			</button>
		</aside>

		<!-- Map -->
		<div bind:this={mapContainer} class="map-container"></div>
	</div>
</div>

<style>
	.map-page {
		display: flex;
		flex-direction: column;
		height: calc(100vh - 53px);
	}

	/* Header */
	.map-header {
		flex-shrink: 0;
		display: flex;
		align-items: center;
		padding: 0.6rem 1.25rem;
		background: var(--surface);
		border-bottom: 1px solid var(--border);
		gap: 1rem;
	}

	.header-left {
		display: flex;
		align-items: baseline;
		gap: 1rem;
	}

	.map-header h1 {
		margin: 0;
		font-size: 1.05rem;
		font-weight: 700;
		color: var(--text);
		letter-spacing: 0.02em;
	}

	.header-subtitle {
		font-size: 0.8rem;
		color: var(--text-muted);
	}

	/* Body */
	.map-body {
		display: flex;
		flex: 1;
		overflow: hidden;
		position: relative;
	}

	/* Sidebar */
	.sidebar {
		position: relative;
		flex-shrink: 0;
		width: 240px;
		background: var(--surface);
		border-right: 1px solid var(--border);
		display: flex;
		overflow: visible;
		transition: width 0.25s ease;
	}

	.sidebar.collapsed {
		width: 0;
	}

	.sidebar-content {
		width: 240px;
		overflow-y: auto;
		overflow-x: hidden;
		flex-shrink: 0;
		padding: 0.75rem 0;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.sidebar.collapsed .sidebar-content {
		visibility: hidden;
	}

	.sidebar-section {
		padding: 0.5rem 1rem;
		border-bottom: 1px solid var(--border);
	}

	.sidebar-section:last-child {
		border-bottom: none;
	}

	.sidebar-section h2 {
		margin: 0 0 0.5rem 0;
		font-size: 0.7rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--text-muted);
	}

	/* Layer switcher (floating, always visible) */
	.layer-switcher {
		position: absolute;
		top: 1rem;
		right: 1rem;
		z-index: 10;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 6px;
		padding: 0.5rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
		display: flex;
		flex-direction: column;
		gap: 0.3rem;
	}

	.layer-switcher-label {
		font-size: 0.68rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--text-muted);
		padding: 0 0.25rem 0.1rem;
	}

	.layer-list {
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.layer-btn {
		background: transparent;
		border: 1px solid transparent;
		border-radius: 4px;
		color: var(--text-muted);
		cursor: pointer;
		font-size: 0.8rem;
		font-family: inherit;
		padding: 0.3rem 0.6rem;
		text-align: left;
		transition: all 0.15s;
		white-space: nowrap;
	}

	.layer-btn:hover {
		color: var(--text);
		background: var(--border);
	}

	.layer-btn.active {
		background: color-mix(in srgb, var(--accent) 15%, transparent);
		color: var(--accent);
		border-color: color-mix(in srgb, var(--accent) 40%, transparent);
	}

	/* Data list */
	.data-list {
		margin: 0;
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 0.25rem 0.75rem;
		font-size: 0.8rem;
	}

	.data-list dt {
		color: var(--text-muted);
	}

	.data-list dd {
		margin: 0;
		color: var(--text);
		text-align: right;
	}

	/* Sidebar toggle tab */
	.sidebar-toggle {
		position: absolute;
		right: -28px;
		top: 50%;
		transform: translateY(-50%);
		width: 28px;
		height: 48px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-left: none;
		border-radius: 0 6px 6px 0;
		color: var(--text-muted);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: color 0.15s, background 0.15s;
		z-index: 10;
	}

	.sidebar-toggle:hover {
		color: var(--text);
		background: var(--border);
	}

	/* Map */
	.map-container {
		flex: 1;
		min-width: 0;
	}
</style>