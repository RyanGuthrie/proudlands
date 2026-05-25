<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import maplibregl, { type StyleSpecification } from 'maplibre-gl';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';

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

	const DEFAULT_CENTER = { lat: 39.9795, lng: -105.5776, zoom: 15.5 };

	let activeLayerId = $state<LayerId>('usgs-topo');
	let sidebarOpen = $state(true);

	let centerLat = $state(DEFAULT_CENTER.lat);
	let centerLng = $state(DEFAULT_CENTER.lng);
	let zoom = $state(6);

	let gotoOpen = $state(false);
	let inputLat = $state('');
	let inputLng = $state('');
	let gotoError = $state('');
	let latInput: HTMLInputElement;

	const COOKIE = 'map_view';
	const COOKIE_MAX_AGE = 60 * 60 * 24 * 365; // 1 year

	function getPositionCookie(): { lat: number; lng: number; zoom: number } | null {
		const entry = document.cookie.split('; ').find((r) => r.startsWith(COOKIE + '='));
		if (!entry) return null;
		try {
			return JSON.parse(decodeURIComponent(entry.slice(COOKIE.length + 1)));
		} catch {
			return null;
		}
	}

	function setPositionCookie(lat: number, lng: number, z: number) {
		const val = encodeURIComponent(JSON.stringify({ lat, lng, zoom: z }));
		document.cookie = `${COOKIE}=${val}; max-age=${COOKIE_MAX_AGE}; path=/; SameSite=Lax`;
	}

	function updateUrl(lat: number, lng: number, z: number) {
		const params = new URLSearchParams(window.location.search);
		params.set('lat', lat.toFixed(6));
		params.set('lng', lng.toFixed(6));
		params.set('zoom', z.toFixed(2));
		history.replaceState(null, '', `${window.location.pathname}?${params}`);
	}

	function getUrlPosition(): { lat: number; lng: number; zoom: number } | null {
		const p = new URLSearchParams(window.location.search);
		const lat = parseFloat(p.get('lat') ?? '');
		const lng = parseFloat(p.get('lng') ?? '');
		const zoom = parseFloat(p.get('zoom') ?? '');
		if (isNaN(lat) || isNaN(lng)) return null;
		return { lat, lng, zoom: isNaN(zoom) ? DEFAULT_CENTER.zoom : zoom };
	}

	function fmt(val: number, decimals: number) {
		return val.toFixed(decimals);
	}

	let coordsCopied = $state(false);

	function copyCoords(lat: number, lng: number) {
		navigator.clipboard.writeText(`${lat.toFixed(6)}, ${lng.toFixed(6)}`);
		coordsCopied = true;
		setTimeout(() => (coordsCopied = false), 2000);
	}

	function toggleGoto() {
		gotoOpen = !gotoOpen;
		gotoError = '';
		if (gotoOpen) {
			// focus after transition starts
			setTimeout(() => latInput?.focus(), 50);
		}
	}

	function gotoCoords() {
		const lat = parseFloat(inputLat);
		const lng = parseFloat(inputLng);
		if (isNaN(lat) || lat < -90 || lat > 90) {
			gotoError = 'Latitude must be between -90 and 90';
			return;
		}
		if (isNaN(lng) || lng < -180 || lng > 180) {
			gotoError = 'Longitude must be between -180 and 180';
			return;
		}
		gotoError = '';
		gotoOpen = false;
		inputLat = '';
		inputLng = '';
		map.flyTo({ center: [lng, lat], zoom: Math.max(zoom, 10), duration: 1200 });
	}

	function handleGotoKey(e: KeyboardEvent) {
		if (e.key === 'Enter') gotoCoords();
		if (e.key === 'Escape') { gotoOpen = false; gotoError = ''; }
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
		const fromUrl = getUrlPosition();
		const saved = fromUrl ?? getPositionCookie();
		const initialCenter: [number, number] = saved ? [saved.lng, saved.lat] : [DEFAULT_CENTER.lng, DEFAULT_CENTER.lat];
		const initialZoom = saved ? saved.zoom : DEFAULT_CENTER.zoom;

		centerLat = initialCenter[1];
		centerLng = initialCenter[0];
		zoom = initialZoom;

		map = new maplibregl.Map({
			container: mapContainer,
			style: makeStyle(activeLayerId),
			center: initialCenter,
			zoom: initialZoom,
		});
		map.addControl(new maplibregl.NavigationControl(), 'bottom-right');
		map.addControl(new maplibregl.ScaleControl(), 'bottom-left');

		map.on('move', () => {
			const c = map.getCenter();
			centerLat = c.lat;
			centerLng = c.lng;
			zoom = map.getZoom();
		});

		map.on('moveend', () => {
			setPositionCookie(centerLat, centerLng, zoom);
			updateUrl(centerLat, centerLng, zoom);
		});

		map.on('style.load', () => {
			applyTrailGeometry(activeTrailGeometry);
		});
	});

	onDestroy(() => {
		map?.remove();
	});

	// Trails
	type TrailDetail = components['schemas']['TrailOutput'];
	type TrailGeometry = components['schemas']['TrailGeometryOutputBody'];

	const SOURCE_ID = 'trail-geometry';
	const LAYER_IDS = ['trail-line-solid', 'trail-line-dashed', 'trail-line-dotted'] as const;

	let activeTrailGeometry = $state<TrailGeometry | null>(null);

	function buildFeatureCollection(geom: TrailGeometry): maplibregl.GeoJSONSourceSpecification['data'] {
		return {
			type: 'FeatureCollection',
			features: (geom.segments ?? []).map((seg) => ({
				type: 'Feature' as const,
				geometry: {
					type: 'LineString' as const,
					coordinates: (seg.coordinates ?? []) as [number, number][],
				},
				properties: { color: seg.color, style: seg.style, width: seg.width },
			})),
		};
	}

	function applyTrailGeometry(geom: TrailGeometry | null) {
		if (!map) return;
		for (const id of LAYER_IDS) {
			if (map.getLayer(id)) map.removeLayer(id);
		}
		if (map.getSource(SOURCE_ID)) map.removeSource(SOURCE_ID);
		if (!geom) return;

		map.addSource(SOURCE_ID, { type: 'geojson', data: buildFeatureCollection(geom) });

		map.addLayer({
			id: 'trail-line-solid',
			type: 'line',
			source: SOURCE_ID,
			filter: ['==', ['get', 'style'], 'solid'],
			layout: { 'line-cap': 'round', 'line-join': 'round' },
			paint: { 'line-color': ['get', 'color'], 'line-width': ['get', 'width'], 'line-opacity': 0.9 },
		});

		map.addLayer({
			id: 'trail-line-dashed',
			type: 'line',
			source: SOURCE_ID,
			filter: ['==', ['get', 'style'], 'dashed'],
			layout: { 'line-cap': 'butt', 'line-join': 'round' },
			paint: { 'line-color': ['get', 'color'], 'line-width': ['get', 'width'], 'line-dasharray': [6, 4], 'line-opacity': 0.9 },
		});

		map.addLayer({
			id: 'trail-line-dotted',
			type: 'line',
			source: SOURCE_ID,
			filter: ['==', ['get', 'style'], 'dotted'],
			layout: { 'line-cap': 'round', 'line-join': 'round' },
			paint: { 'line-color': ['get', 'color'], 'line-width': ['get', 'width'], 'line-dasharray': [1, 4], 'line-opacity': 0.9 },
		});
	}

	const TRAIL_PAGE_SIZE = 5;

	let trailNames = $state<string[]>([]);
	let trailsError = $state('');
	let trailSearch = $state('');
	let showAllTrails = $state(false);
	let selectedTrail = $state('');
	let trailDetail = $state<TrailDetail | null>(null);
	let trailDetailError = $state('');
	let trailDetailLoading = $state(false);

	let filteredTrails = $derived(
		[...(trailSearch.trim()
			? trailNames.filter((n) => n.toLowerCase().includes(trailSearch.toLowerCase().trim()))
			: trailNames)
		].sort((a, b) => a.localeCompare(b))
	);

	let displayedTrails = $derived(
		showAllTrails || trailSearch.trim() ? filteredTrails : filteredTrails.slice(0, TRAIL_PAGE_SIZE)
	);

	let hiddenCount = $derived(
		!trailSearch.trim() && !showAllTrails ? Math.max(0, filteredTrails.length - TRAIL_PAGE_SIZE) : 0
	);

	async function fetchTrails() {
		const { data, error } = await api.GET('/trail');
		if (error || !data) {
			trailsError = 'Could not load trails';
			return;
		}
		trailNames = data.resources ?? [];
	}

	async function selectTrail(name: string) {
		if (selectedTrail === name) {
			selectedTrail = '';
			trailDetail = null;
			activeTrailGeometry = null;
			applyTrailGeometry(null);
			return;
		}
		selectedTrail = name;
		trailDetail = null;
		trailDetailError = '';
		trailDetailLoading = true;
		activeTrailGeometry = null;
		applyTrailGeometry(null);

		const [metaResult, geoResult] = await Promise.all([
			api.GET('/trail/{name}', { params: { path: { name } } }),
			api.GET('/trail/{name}/geometry', { params: { path: { name } } }),
		]);

		trailDetailLoading = false;

		if (metaResult.error || !metaResult.data) {
			trailDetailError = 'Could not load trail details';
		} else {
			trailDetail = metaResult.data;
		}

		if (!geoResult.error && geoResult.data) {
			activeTrailGeometry = geoResult.data;
			applyTrailGeometry(activeTrailGeometry);
		}
	}

	async function flyToTrail(name: string) {
		let detail: TrailDetail | null = selectedTrail === name ? trailDetail : null;
		if (!detail) {
			const { data } = await api.GET('/trail/{name}', { params: { path: { name } } });
			if (data) detail = data;
		}
		if (detail) {
			map.flyTo({ center: [detail.longitude, detail.latitude], zoom: Math.max(zoom, 12), duration: 1200 });
		}
	}

	onMount(() => {
		fetchTrails();
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

				<section class="sidebar-section trails-section">
					<h2>Trails</h2>

					<div class="trail-search-wrap">
						<svg class="trail-search-icon" width="12" height="12" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
							<circle cx="6.5" cy="6.5" r="5"/>
							<line x1="10.5" y1="10.5" x2="15" y2="15"/>
						</svg>
						<input
							class="trail-search"
							type="search"
							placeholder="Search trails…"
							bind:value={trailSearch}
						/>
					</div>

					{#if trailsError}
						<p class="trail-error">{trailsError}</p>
					{:else if trailNames.length === 0}
						<p class="trail-empty">No trails found</p>
					{:else if filteredTrails.length === 0}
						<p class="trail-empty">No trails match "{trailSearch}"</p>
					{:else}
						<ul class="trail-list">
							{#each displayedTrails as name}
								<li class="trail-row">
									<button
										class="trail-item"
										class:active={selectedTrail === name}
										onclick={() => selectTrail(name)}
									>
										{name}
									</button>
									<button
										class="trail-fly"
										onclick={() => flyToTrail(name)}
										title="Fly to trail"
									>
										<svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
											<path d="M8 0a5 5 0 0 0-5 5c0 4 5 11 5 11s5-7 5-11a5 5 0 0 0-5-5zm0 7a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"/>
										</svg>
									</button>
								</li>
							{/each}
						</ul>
						{#if hiddenCount > 0}
							<button class="trail-more" onclick={() => (showAllTrails = true)}>
								+{hiddenCount} more
							</button>
						{:else if showAllTrails && filteredTrails.length > TRAIL_PAGE_SIZE}
							<button class="trail-more" onclick={() => (showAllTrails = false)}>
								Show less
							</button>
						{/if}
					{/if}

					{#if selectedTrail}
						<div class="trail-detail">
							{#if trailDetailLoading}
								<p class="trail-loading">Loading…</p>
							{:else if trailDetailError}
								<p class="trail-error">{trailDetailError}</p>
							{:else if trailDetail}
								<div class="trail-detail-header">
									<span class="trail-name">{trailDetail.name}</span>
									<span class="trail-difficulty difficulty-{trailDetail.difficulty}">{trailDetail.difficulty}</span>
								</div>
								<p class="trail-description">{trailDetail.description}</p>
								<dl class="data-list">
									<dt>Length</dt>
									<dd>{trailDetail.length_miles} mi</dd>
									<dt>Latitude</dt>
									<dd>{fmt(Math.abs(trailDetail.latitude), 4)}° {trailDetail.latitude >= 0 ? 'N' : 'S'}</dd>
									<dt>Longitude</dt>
									<dd>{fmt(Math.abs(trailDetail.longitude), 4)}° {trailDetail.longitude >= 0 ? 'E' : 'W'}</dd>
								</dl>
								<div class="location-actions">
									<button class="loc-btn" onclick={() => copyCoords(trailDetail!.latitude, trailDetail!.longitude)}>
										{#if coordsCopied}
											<svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor"><path d="M13.5 2.5L6 10 2.5 6.5 1 8l5 5 9-9z"/></svg>
											Copied
										{:else}
											<svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor"><path d="M4 4h8v10H4zM2 2h8v1H2zM6 0h8v12h-1V1H6z"/></svg>
											Copy coords
										{/if}
									</button>
									<button class="loc-btn" onclick={() => flyToTrail(selectedTrail)}>
										<svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor"><path d="M8 0a5 5 0 0 0-5 5c0 4 5 11 5 11s5-7 5-11a5 5 0 0 0-5-5zm0 7a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"/></svg>
										Go to
									</button>
								</div>
							{/if}
						</div>
					{/if}
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

		<!-- Map area: overlays are scoped here so they stay within the map -->
		<div class="map-area">
			<!-- Coordinate jump popout -->
			<div class="goto-wrap" class:open={gotoOpen}>
				<button class="goto-trigger" onclick={toggleGoto} aria-label="Go to coordinates" title="Go to coordinates">
					<svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
						<circle cx="8" cy="8" r="3"/>
						<line x1="8" y1="1" x2="8" y2="4"/>
						<line x1="8" y1="12" x2="8" y2="15"/>
						<line x1="1" y1="8" x2="4" y2="8"/>
						<line x1="12" y1="8" x2="15" y2="8"/>
					</svg>
				</button>
				<div class="goto-form">
					<input
						bind:this={latInput}
						type="text"
						inputmode="decimal"
						placeholder="Latitude"
						bind:value={inputLat}
						onkeydown={handleGotoKey}
					/>
					<span class="goto-sep">,</span>
					<input
						type="text"
						inputmode="decimal"
						placeholder="Longitude"
						bind:value={inputLng}
						onkeydown={handleGotoKey}
					/>
					<button class="goto-go" onclick={gotoCoords}>Go</button>
				</div>
				{#if gotoError}
					<div class="goto-error">{gotoError}</div>
				{/if}
			</div>

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

			<div bind:this={mapContainer} class="map-container"></div>
		</div>
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
	}

	.map-area {
		flex: 1;
		position: relative;
		min-width: 0;
		overflow: hidden;
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
		width: 100%;
		height: 100%;
	}

	/* Coordinate jump popout */
	.goto-wrap {
		position: absolute;
		top: 1rem;
		left: 1rem;
		z-index: 10;
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.4rem;
	}

	.goto-trigger {
		width: 36px;
		height: 36px;
		border-radius: 6px;
		background: var(--surface);
		border: 1px solid var(--border);
		color: var(--text-muted);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
		transition: color 0.15s, background 0.15s;
		flex-shrink: 0;
	}

	.goto-trigger:hover,
	.goto-wrap.open .goto-trigger {
		color: var(--accent);
		background: color-mix(in srgb, var(--accent) 10%, var(--surface));
		border-color: color-mix(in srgb, var(--accent) 40%, transparent);
	}

	.goto-form {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 6px;
		padding: 0.4rem 0.5rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
		max-width: 0;
		overflow: hidden;
		opacity: 0;
		transition: max-width 0.25s ease, opacity 0.2s ease, padding 0.25s ease;
		padding-left: 0;
		padding-right: 0;
		white-space: nowrap;
	}

	.goto-wrap.open .goto-form {
		max-width: 340px;
		opacity: 1;
		padding: 0.4rem 0.5rem;
	}

	.goto-form input {
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: 4px;
		color: var(--text);
		font-family: inherit;
		font-size: 0.82rem;
		padding: 0.3rem 0.5rem;
		width: 120px;
		outline: none;
		transition: border-color 0.15s;
	}

	.goto-form input:focus {
		border-color: var(--accent);
	}

	.goto-form input::placeholder {
		color: var(--text-muted);
	}

	.goto-sep {
		color: var(--text-muted);
		font-size: 0.9rem;
		flex-shrink: 0;
	}

	.goto-go {
		background: var(--accent);
		border: none;
		border-radius: 4px;
		color: #0f1117;
		cursor: pointer;
		font-family: inherit;
		font-size: 0.8rem;
		font-weight: 600;
		padding: 0.3rem 0.7rem;
		flex-shrink: 0;
		transition: background 0.15s;
	}

	.goto-go:hover {
		background: var(--accent-hover);
	}

	.goto-error {
		background: color-mix(in srgb, #e05c5c 12%, var(--surface));
		border: 1px solid color-mix(in srgb, #e05c5c 40%, transparent);
		border-radius: 4px;
		color: #f08080;
		font-size: 0.75rem;
		padding: 0.3rem 0.6rem;
		white-space: nowrap;
	}

	/* Location actions */
	.location-actions {
		display: flex;
		gap: 0.4rem;
		margin-top: 0.5rem;
	}

	.loc-btn {
		display: flex;
		align-items: center;
		gap: 0.3rem;
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: 4px;
		color: var(--text-muted);
		cursor: pointer;
		font-family: inherit;
		font-size: 0.75rem;
		padding: 0.25rem 0.5rem;
		transition: color 0.15s, border-color 0.15s;
	}

	.loc-btn:hover {
		color: var(--text);
		border-color: var(--text-muted);
	}

	/* Trails */
	.trail-search-wrap {
		position: relative;
		display: flex;
		align-items: center;
	}

	.trail-search-icon {
		position: absolute;
		left: 0.45rem;
		color: var(--text-muted);
		pointer-events: none;
		flex-shrink: 0;
	}

	.trail-search {
		width: 100%;
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: 4px;
		color: var(--text);
		font-family: inherit;
		font-size: 0.78rem;
		padding: 0.3rem 0.5rem 0.3rem 1.75rem;
		outline: none;
		transition: border-color 0.15s;
	}

	.trail-search:focus {
		border-color: var(--accent);
	}

	.trail-search::placeholder {
		color: var(--text-muted);
	}

	/* hide the native clear button on search inputs */
	.trail-search::-webkit-search-cancel-button {
		-webkit-appearance: none;
	}

	.trail-more {
		background: none;
		border: none;
		color: var(--accent);
		cursor: pointer;
		font-family: inherit;
		font-size: 0.78rem;
		padding: 0.1rem 0;
		text-align: left;
		transition: color 0.15s;
	}

	.trail-more:hover {
		color: var(--accent-hover);
	}

	.trails-section {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.trail-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.trail-row {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.trail-item {
		background: transparent;
		border: 1px solid transparent;
		border-radius: 4px;
		color: var(--text-muted);
		cursor: pointer;
		font-family: inherit;
		font-size: 0.82rem;
		padding: 0.3rem 0.5rem;
		text-align: left;
		transition: all 0.15s;
		flex: 1;
		min-width: 0;
		text-transform: capitalize;
	}

	.trail-item:hover {
		color: var(--text);
		background: var(--border);
	}

	.trail-item.active {
		background: color-mix(in srgb, var(--accent) 15%, transparent);
		color: var(--accent);
		border-color: color-mix(in srgb, var(--accent) 40%, transparent);
	}

	.trail-fly {
		background: transparent;
		border: none;
		border-radius: 4px;
		color: var(--text-muted);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		padding: 0.3rem;
		transition: color 0.15s, background 0.15s;
	}

	.trail-fly:hover {
		color: var(--accent);
		background: color-mix(in srgb, var(--accent) 10%, transparent);
	}

	.trail-detail {
		border-top: 1px solid var(--border);
		padding-top: 0.6rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.trail-detail-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
	}

	.trail-name {
		font-size: 0.85rem;
		font-weight: 600;
		color: var(--text);
	}

	.trail-difficulty {
		font-size: 0.68rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.06em;
		padding: 0.15rem 0.4rem;
		border-radius: 3px;
		flex-shrink: 0;
	}

	.difficulty-easy {
		background: color-mix(in srgb, #4caf50 15%, transparent);
		color: #81c784;
	}

	.difficulty-moderate {
		background: color-mix(in srgb, #ff9800 15%, transparent);
		color: #ffb74d;
	}

	.difficulty-hard {
		background: color-mix(in srgb, #f44336 15%, transparent);
		color: #e57373;
	}

	.trail-description {
		margin: 0;
		font-size: 0.8rem;
		color: var(--text-muted);
		line-height: 1.5;
	}

	.trail-empty,
	.trail-loading,
	.trail-error {
		margin: 0;
		font-size: 0.8rem;
		color: var(--text-muted);
	}

	.trail-error {
		color: #f08080;
	}
</style>