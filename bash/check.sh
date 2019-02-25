#!/usr/bin/env bash
check(){
  echo "Checking the availablility of the keytab"
  kadmin -kt /tmp/admin.keytab -p admin/admin -q "listprincs"
}