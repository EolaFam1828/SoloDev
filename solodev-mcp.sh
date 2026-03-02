#!/bin/bash
# Wrapper script for SoloDev MCP server.
# Ensures the working directory is correct so the SQLite database can be found.
cd /Users/mjc01/SoloDev

if [ -x ./solodev ]; then
  exec ./solodev mcp stdio .local.env
fi

exec ./gitness mcp stdio .local.env
