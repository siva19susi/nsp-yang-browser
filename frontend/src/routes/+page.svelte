<script lang="ts">
	import Footer from '$lib/components/Footer.svelte'

</script>
  
<svelte:head>
	<title>NSP YANG Browser</title>
</svelte:head>

<div class="flex flex-col items-center min-h-screen pt-5 has-header-img font-nunito">
  <div class="flex-grow-1 flex-shrink-0 m-auto px-4 py-10">
    <p class="mb-4"><img src="/images/nblue.png" width="100" alt="Logo" /></p>
    <h3 class="text-3xl font-light mt-6">NSP YANG Browser</h3>
    <div class="text-white mt-4 max-w-[550px]">
      <p>
        NSP allows you to create and execute intent-based automation model for implementing network-level 
        changes that translates high-level designs into necessary network configuration. In order to 
        faciliate this flow, YANG modeling is used to leverage the automation to its complete heights.
      </p>
      <p class="mt-4">
        This tool provides a means to explore these YANG models either from a connected NSP system 
        or the YANG repository can be manually uploaded and explored.
      </p>
    </div>
    <div class="grid grid-cols-2 gap-3 pt-10">
      <a class="px-2 py-2 rounded-full text-nowrap text-center text-white bg-gray-600 hover:bg-gray-700" href="/nsp">NSP Connect</a>
      <a class="px-2 py-2 rounded-full text-nowrap text-center text-white bg-gray-600 hover:bg-gray-700" href="/uploads">Uploads</a>
    </div>
  </div>
  <div class="flex-grow-0 flex-shrink-0">
    <Footer home={true} />
  </div>
</div>