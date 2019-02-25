#!/usr/bin/env bash
keytab(){
  echo "Generating keytab"
  kadmin -kt /tmp/admin.keytab -p admin/admin -q "addprinc -randkey $1"
  kadmin -kt /tmp/admin.keytab -p admin/admin -q "ktadd -k $2 $1"
}