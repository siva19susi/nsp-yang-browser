<script lang="ts">
  export let x: string
  export let y: string
</script>

<div class="flex h-screen p-7 font-nunito">
  <div class="m-auto">
    <article class="rounded max-w-[550px]">
			<div class="text-white bg-gray-500 px-5 py-3 rounded-t flex items-center">
        <a href="../" class="flex mr-4"><img src="/images/navbar-logo.png" alt="Logo" width="25"/></a>
				<div>
          <p class="text-xs">Yang Browser</p>
          <p>Nothing to Compare</p>
        </div>
			</div>
			<div class="bg-gray-200 p-5 rounded-b">
        <p>No path or path types have changed between {x} and {y}. </p>
        <p class="mt-6">However, there might be changes in descriptions, type constraints, or other yang statements. 
          To identify these changes, kindly contact the repo source Admin.</p>
			</div>
		</article>
  </div>
</div>