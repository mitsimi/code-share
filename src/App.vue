<template>
  <RouterView />
</template>

<style>
/* Custom scrollbar styling */
::-webkit-scrollbar {
  width: 16px;
  height: 16px;
}

::-webkit-scrollbar-track {
  background: #fcf7de;
  border-left: 3px solid #000000;
}

::-webkit-scrollbar-thumb {
  background: #ff3e4d;
  border: 3px solid #000000;
  box-shadow: inset 2px 2px 0 rgba(255, 255, 255, 0.3);
}

::-webkit-scrollbar-corner {
  background: #fcf7de;
}
</style>
