let timer: any = undefined;

export function push() {
  clearTimeout(timer);
  timer = null;
  window.history.pushState(
    window.history.state,
    "",
    `${window.location.pathname}${window.location.search}#${window.location.hash}`,
  );
}
export function throttledPushState(searchParams: string) {
  if (!timer) {
    timer = setTimeout(push, 60000);
  }
  location.replace(`#${searchParams}`);
}
