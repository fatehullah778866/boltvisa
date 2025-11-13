<!-- 47a443b1-7464-4e6e-9511-ec75bf6cc6c4 66a7fc43-83f9-4001-911a-b9ee33e64c28 -->
# Fix Health Check and Request Timeout Issues

## Problems Identified

1. Health check is aborting with "signal is aborted without reason" - likely timing out or proxy not working
2. Login request timing out after 30 seconds - suggests backend unreachable or proxy misconfiguration
3. Health check is causing more problems than it solves

## Solution

1. **Remove the health check entirely** - It's non-essential and causing issues. The actual request will provide better error messages.
2. **Ensure API requests use Next.js proxy correctly** - Verify `/api` paths are properly proxied
3. **Fix URL construction** - Make sure relative paths use Next.js proxy, absolute URLs work directly
4. **Improve error handling** - Better timeout handling and error messages
5. **Add request logging** - Help debug connection issues

## Changes

### File: `apps/web/src/lib/apiRequest.ts`

1. **Remove health check function and call** (lines 35-76, 98-110):

- Remove `checkBackendHealth` function entirely
- Remove health check call from `apiRequest`
- This eliminates the abort errors

2. **Fix URL construction** (lines 90-95):

- Ensure `/api` paths are used as-is (Next.js will proxy them)
- Verify absolute URLs work correctly
- Add logging to see what URL is being used

3. **Improve timeout handling**:

- Reduce timeout from 30s to 10s for faster feedback
- Better error messages for timeouts
- Distinguish between timeout and connection errors

4. **Add request logging** (optional, for debugging):

- Log the URL being requested
- Log timeout/error details

## Implementation Details

- Remove health check completely - it's redundant and problematic
- Use Next.js proxy for all `/api` paths (already configured in next.config.js)
- Direct URLs only for external APIs if needed
- Better error messages that help users understand what's wrong

### To-dos

- [ ] Remove health check function and all references to it
- [ ] Fix URL construction to ensure Next.js proxy is used correctly
- [ ] Improve timeout handling and error messages
- [ ] Test that requests work correctly with Next.js proxy