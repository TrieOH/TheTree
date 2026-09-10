# Issue draft — CachyOS/distribution

> Copy the body below into https://github.com/CachyOS/distribution/issues/new
> Before submitting: run `sudo cachyos-bugreport.sh` and paste the paste.cachyos.org link into the "Bug report" section.

---

## Title

sudo/su/passwd reject the correct password on the stock pambase auth chain — `PAM_BAD_JUMP` from `auth [success=1 default=bad] pam_unix.so`; changing the control to `required` fixes it (recurring since ~Sep 2025, not faillock)

## Environment

- Distro: CachyOS (rolling), kernel `7.1.5-1-cachyos`
- DE: KDE Plasma 6 (Wayland), greeter `plasmalogin`
- pam `1.7.2-2.1` · pambase `20260616-1` · sudo `1.9.17.p2-6.1` · systemd `261.2-1` · libxcrypt `4.5.2-1.1` · glibc `2.44+r3+g0b05bc142249-1`
- Account: local user in `wheel` and `nopasswdlogin`; fish shell

## Summary

`sudo` (and `su`, and the current-password check in `passwd`) reject the **correct** password. The password was verified byte-for-byte against the `/etc/shadow` hash (`$y$`/yescrypt, unchanged since install) via the host's `libcrypt.so.2` (the library `pam_unix.so` actually links). Logs never show a classic wrong-password `authentication failure` — instead every attempt logs:

```
pam_unix(sudo:auth): conversation failed
pam_unix(sudo:auth): auth could not identify password for [user]
user : N incorrect password attempts ; ... COMMAND=...
```

`auth could not identify password` is the pam_unix message when `pam_get_authtok()` fails (linux-pam `pam_unix_auth.c`), i.e. the token never reaches the hash comparison. The greeter login keeps working, which is why it looks like "sudo is broken but my password is fine" — it is fine.

Root cause found by bisecting the auth chain: the stock control flag **`[success=1 default=bad]`** on the `pam_unix.so` line makes libpam fail (`PAM_BAD_JUMP`, "bad jump in stack") when `pam_unix` returns success on this setup. Replacing it with **`required`** (functionally equivalent for a single auth module) fixes auth completely.

## Repro (deterministic)

Setup: stock `/etc/pam.d/sudo` = `auth include system-auth`; stock `/etc/pam.d/system-auth` auth section:

```
-auth      [success=2 default=ignore]  pam_systemd_home.so
auth       [success=1 default=bad]     pam_unix.so          try_first_pass nullok
auth       optional                    pam_permit.so
auth       required                    pam_env.so
```

1. `printf '%s\n' '<correct password>' | sudo -k -S true` → **fails** (rc 1, ~2 s delay = pam_unix's `pam_fail_delay` 2000000 µs), journal shows the pair above.
2. Same command with a deliberately wrong password → also fails (expected), but journal shows classic `authentication failure` — proving the correct password is never even compared in case 1.
3. Same account, same password, same terminal, same session — replace only the `pam_unix.so` auth line's control with `required`:

```
auth       required                   pam_unix.so          try_first_pass nullok
```

→ same command **succeeds** (rc 0, ~60 ms). Wrong password is still rejected.

Bisection matrix (3 runs each, same password, only the auth section of a test service changed):

| auth chain (pam_unix args identical) | result |
|---|---|
| `required pam_unix.so` | OK (~60 ms) |
| `[success=1 default=bad] pam_unix.so` (only module) | FAIL, `bad jump in stack` |
| `[success=1 default=bad] pam_unix.so` + `optional pam_permit` + `required pam_env` | FAIL (conv errors) |
| `[success=1 default=bad] pam_unix.so` + pam_systemd_home before (stock, incl. via include) | FAIL |
| `required pam_unix.so` + permit + env | OK |

`journalctl` shows `PAM bad jump in stack` (libpam `pam_dispatch.c`: a `success=N` control whose skip cannot be satisfied is treated as a config syntax error → fail).

## Not faillock

faillock is **commented out** in this pambase's system-auth (since pambase 20260616-1), `faillock` shows nothing, and a reboot does **not** fix it (most forum reports say reboot fixes it; here it persisted across a fresh boot and a re-login — fully deterministic).

## Context observations (may or may not be load-bearing)

- `pam_systemd_home.so` (shipped by `systemd`, 261.2-1) runs first in the chain; on this machine it tries D-Bus activation of `org.freedesktop.home1` and fails: `Could not activate remote peer 'org.freedesktop.home1': activation request failed: unknown unit` (systemd-homed installed but disabled). Removing the homed line from the test chain does **not** fix the failure, so it is not the (sole) trigger.
- Account is in `nopasswdlogin` (KDE greeter can log in without verifying the password against shadow — which is why "login works" is a red herring here).
- Workaround shipped on the affected machine: system-auth line changed to `required` (backup kept). Fine to keep until root-caused, but the stock config should not behave this way.

## Links to existing reports (same log signature, no root cause yet)

- https://discuss.cachyos.org/t/password-the-correct-one-doesnt-work-randomly-after-fresh-install-after-an-update-and-even-when-witthin-a-logged-in-session-in-terminal-and-frequently/15747 (open since Sep 2025, 34 posts, last Aug 2026)
- https://discuss.cachyos.org/t/sudo-password-stops-working-suddenly/24377
- https://discuss.cachyos.org/t/sudo-password-failure-loop-despite-correct-credentials-and-working-root-login/30370
- https://discuss.cachyos.org/t/cant-install-updates-password-is-always-incorrect-even-though-its-correct/26855

## Bug report

<paste link from `sudo cachyos-bugreport.sh` here>

## Questions for maintainers

1. Does the stock chain reproduce on a vanilla Arch install, or is it specific to the CachyOS `pam 1.7.2-2.1` rebuild vs Arch `1.7.2-2` (or the `systemd 261.2` pam_systemd_home)?
2. Is `[success=1 default=bad]` + `try_first_pass nullok` + a successful pam_unix known to be problematic with linux-pam 1.7.2's substack/include dispatch (libpam `pam_dispatch.c`)? If yes, this should be escalated to linux-pam upstream.
3. What changed around the recent boots that makes this latent config flip from working to failing with **zero package/config changes** in between?
