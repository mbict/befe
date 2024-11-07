#!/bin/sh
#this script regenerates the dist.go file that is used to determine the entry points for each project

WORKDIR=$(pwd)

# copy partial template
cp $WORKDIR/dist.go.template $WORKDIR/dist.go

# list the index entry points for all the projects
cd $WORKDIR/dist/assets/account
echo '\t"account": { JS:"'$(ls -1 -p index-*.js)'", CSS:"'$(ls -1 -p index-*.css)'" },' >> $WORKDIR/dist.go
cd $WORKDIR/dist/assets/backoffice
echo '\t"backoffice": { JS:"'$(ls -1 -p index-*.js)'", CSS:"'$(ls -1 -p index-*.css)'" },' >> $WORKDIR/dist.go
cd $WORKDIR/dist/assets/agenda
echo '\t"agenda": { JS:"'$(ls -1 -p index-*.js)'", CSS:"'$(ls -1 -p index-*.css)'" },' >> $WORKDIR/dist.go
cd $WORKDIR/dist/assets/admin
echo '\t"admin": { JS:"'$(ls -1 -p index-*.js)'", CSS:"'$(ls -1 -p index-*.css)'" },' >> $WORKDIR/dist.go

echo "}" >> $WORKDIR/dist.go
