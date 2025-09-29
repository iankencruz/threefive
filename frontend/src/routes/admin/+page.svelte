<script lang="ts">
// Will add dashboard logic here later
</script>

<div class="dashboard">
	<h1>Admin Dashboard</h1>
	
	<div class="stats">
		<div class="stat-card">
			<h3>Total Users</h3>
			<p class="stat-number">0</p>
		</div>
		
		<div class="stat-card">
			<h3>Active Sessions</h3>
			<p class="stat-number">0</p>
		</div>
		
		<div class="stat-card">
			<h3>New Today</h3>
			<p class="stat-number">0</p>
		</div>
	</div>
	
	<div class="recent-activity">
		<h2>Recent Activity</h2>
		<p>Activity feed will be added here</p>
	</div>
</div>

<style>
	.dashboard h1 {
		margin-bottom: 2rem;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1.5rem;
		margin-bottom: 3rem;
	}

	.stat-card {
		background: white;
		padding: 1.5rem;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.stat-card h3 {
		font-size: 0.875rem;
		color: #666;
		margin-bottom: 0.5rem;
		font-weight: 500;
	}

	.stat-number {
		font-size: 2rem;
		font-weight: bold;
		color: #333;
		margin: 0;
	}

	.recent-activity {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.recent-activity h2 {
		margin-bottom: 1rem;
		font-size: 1.25rem;
	}
</style>
