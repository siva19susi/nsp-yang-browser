<script lang="ts">
	import Footer from '$lib/components/Footer.svelte'

</script>
  
<svelte:head>
	<title>NSP YANG Browser</title>
</svelte:head>

<div class="flex flex-col min-h-screen has-header-img font-nunito">
  <div class="flex grow items-center justify-start px-8 md:px-14 py-10">
    <div>
      <p class="mb-4"><img src="/images/nwhite.png" width="95" alt="Logo" /></p>
      <h3 class="text-3xl text-yellow-300 font-light mt-6">NSP YANG Browser</h3>
      <div class="text-gray-300 mt-4 max-w-[500px] text-sm">
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
      <div class="grid grid-cols-2 gap-3 mt-6 text-sm">
        <a class="px-2 py-2 rounded-lg text-nowrap text-center text-white bg-gray-600 hover:bg-gray-700" href="/nsp">NSP Connect</a>
        <!--<a class="px-2 py-2 rounded-lg text-nowrap text-center text-white bg-gray-600 hover:bg-gray-700" href="/uploads">Offline</a>-->
      </div>
    </div>
  </div>
  <div class="shrink-0">
    <Footer home={true} />
  </div>
</div>