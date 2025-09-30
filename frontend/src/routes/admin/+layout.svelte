<script lang="ts">
// Will add auth check here later
const js = "js"
console.log("js", js)
</script>

<div class="admin-layout">
	<aside class="sidebar">
		<h2>Admin</h2>
		<nav>
			<a href="/admin">Dashboard</a>
			<a href="/admin/users">Users</a>
			<a href="/admin/settings">Settings</a>
		</nav>
	</aside>
	
	<main class="content">
		<slot />
	</main>
</div>

<style>
	.admin-layout {
		display: flex;
		min-height: 100vh;
	}

	.sidebar {
		width: 250px;
		background: #f5f5f5;
		padding: 2rem 1rem;
		border-right: 1px solid #ddd;
	}

	.sidebar h2 {
		margin-bottom: 2rem;
	}

	.sidebar nav {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.sidebar a {
		padding: 0.5rem 1rem;
		text-decoration: none;
		color: #333;
		border-radius: 4px;
	}

	.sidebar a:hover {
		background: #e0e0e0;
	}

	.content {
		flex: 1;
		padding: 2rem;
	}
</style>
