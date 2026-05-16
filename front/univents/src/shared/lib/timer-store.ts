import { useSyncExternalStore } from "react"

let now = Date.now()
const listeners = new Set<() => void>()

let interval: ReturnType<typeof setInterval> | null = null

function start() {
  if (interval !== null) return

  interval = setInterval(() => {
    now = Date.now()
    listeners.forEach((l) => { l() })
  }, 1000)
}

function stop() {
  if (interval === null) return

  clearInterval(interval)
  interval = null
}

function subscribe(callback: () => void) {
  listeners.add(callback)

  // start clock if first subscriber
  if (listeners.size === 1) start()


  return () => {
    listeners.delete(callback)
    // stop clock if nobody listening
    if (listeners.size === 0) stop()
  }
}

function getSnapshot() {
  return now
}

export function useNow() {
  return useSyncExternalStore(subscribe, getSnapshot)
}