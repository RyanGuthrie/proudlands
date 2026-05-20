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
			tileSize: 256,
			attribution: 'USGS National Map',
		},
		{
			id: 'usgs-imagery',
			label: 'USGS Imagery',
			tiles: ['https://basemap.nationalmap.gov/arcgis/rest/services/USGSImageryOnly/MapServer/tile/{z}/{y}/{x}'],
			tileSize: 256,
			attribution: 'USGS National Map',
		},
		{
			id: 'osm',
			label: 'OpenStreetMap',
			tiles: ['https://tile.openstreetmap.org/{z}/{x}/{y}.png'],
			tileSize: 256,
			attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
		},
		{
			id: 'stamen-terrain',
			label: 'Stamen Terrain',
			tiles: ['https://tiles.stadiamaps.com/tiles/stamen_terrain/{z}/{x}/{y}.png'],
			tileSize: 256,
			attribution: '© <a href="https://stadiamaps.com/">Stadia Maps</a> © <a href="https://stamen.com/">Stamen Design</a> © <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
		},
	] as const;

	type LayerId = (typeof layers)[number]['id'];

	let activeLayerId = $state<LayerId>('usgs-topo');

	function makeStyle(layerId: LayerId): StyleSpecification {
		const layer = layers.find((l) => l.id === layerId)!;
		return {
			version: 8,
			sources: {
				tiles: {
					type: 'raster',
					tiles: [...layer.tiles],
					tileSize: layer.tileSize,
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
	});

	onDestroy(() => {
		map?.remove();
	});
</script>

<div class="map-page">
	<div class="layer-switcher">
		{#each layers as layer}
			<button
				class:active={activeLayerId === layer.id}
				onclick={() => switchLayer(layer.id)}
			>
				{layer.label}
			</button>
		{/each}
	</div>
	<div bind:this={mapContainer} class="map-container"></div>
</div>

<style>
	.map-page {
		position: relative;
		height: calc(100vh - 53px);
		width: 100%;
		overflow: hidden;
	}

	.map-container {
		width: 100%;
		height: 100%;
	}

	.layer-switcher {
		position: absolute;
		top: 1rem;
		left: 1rem;
		z-index: 10;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 6px;
		padding: 0.4rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
	}

	button {
		background: transparent;
		border: 1px solid transparent;
		border-radius: 4px;
		color: var(--text-muted);
		cursor: pointer;
		font-size: 0.8rem;
		font-family: inherit;
		padding: 0.35rem 0.75rem;
		text-align: left;
		transition: all 0.15s;
		white-space: nowrap;
	}

	button:hover {
		color: var(--text);
		background: var(--border);
	}

	button.active {
		background: var(--accent);
		color: #0f1117;
		font-weight: 600;
	}
</style>