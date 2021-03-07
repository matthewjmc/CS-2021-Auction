#!/bin/bash
echo "Number of Connection:"
netstat -ant | grep ESTABLISHED | wc -l
