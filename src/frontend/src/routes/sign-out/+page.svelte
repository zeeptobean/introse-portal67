<script lang="ts">
  import { enhance } from "$app/forms";
</script>

<main class="sign-in">
  <h1>Sign out</h1>
  <form method="POST" use:enhance><button class="primary">Sign out</button></form>
  <a href="/account">Return to account</a>
</main>
