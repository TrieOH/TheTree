export function randomItem<T>(items: T[]): T | undefined {
  if (items.length === 0) return undefined;
  const value = new Uint32Array(1);
  crypto.getRandomValues(value);
  return items[value[0] % items.length];
}

export function drawSequence<T>(items: T[], minimumFrames = 22): T[] {
  if (items.length === 0) return [];
  const targetLength = Math.max(minimumFrames, Math.ceil(items.length * 1.5));
  const sequence: T[] = [];
  while (sequence.length < targetLength) {
    const round = shuffle(items);
    const previous = sequence.at(-1);
    if (round.length > 1 && round[0] === previous) {
      const swapIndex = round.findIndex((item) => item !== previous);
      [round[0], round[swapIndex]] = [round[swapIndex], round[0]];
    }
    sequence.push(...round.slice(0, targetLength - sequence.length));
  }
  return sequence;
}

function shuffle<T>(items: T[]): T[] {
  const shuffled = [...items];
  for (let index = shuffled.length - 1; index > 0; index -= 1) {
    const target = new Uint32Array(1);
    crypto.getRandomValues(target);
    const swapIndex = target[0] % (index + 1);
    [shuffled[index], shuffled[swapIndex]] = [
      shuffled[swapIndex],
      shuffled[index],
    ];
  }
  return shuffled;
}

export function drawTimeline(frameCount: number, targetDurationMs = 4800) {
  if (frameCount === 0) return { delays: [], durationMs: 0 };
  const weights = Array.from({ length: frameCount }, (_, index) => {
    const progress = frameCount === 1 ? 1 : index / (frameCount - 1);
    return 0.45 + 1.55 * progress ** 2;
  });
  const scale =
    targetDurationMs / weights.reduce((sum, weight) => sum + weight, 0);
  let elapsed = 0;
  const delays = weights.map((weight) => {
    const delay = elapsed;
    elapsed += Math.max(16, weight * scale);
    return delay;
  });
  return { delays, durationMs: elapsed };
}
