############################################################
#  brief:cp remote file to local
#  date :2020-01-16
############################################################
function dataop {
  dt="$1"
  echo "scp -r /home/stock/stockdata/${dt}"
  #tar
  echo "tar ${dt} file"
  ssh stock@47.89.240.84 "cd /home/stock/stockdata/ && tar -cvf ${curdate}.tar ${curdate}"

  #cp
  echo "scp -r /home/stock/stockdata/${dt}"
  scp -r stock@47.89.240.84:/home/stock/stockdata/${dt}.tar /home/stock/stockdata/

  #
  echo "unzip tar file"
  cd /home/stock/stockdata/
  tar -xvf ${dt}.tar

  #delete local tar file
  rm -f /home/stock/stockdata/${dt}.tar

  #delete remote tar file
  ssh stock@47.89.240.84 "rm -f /home/stock/stockdata/${dt}.tar"

}


startdate="20190709"
while [ "$curdate" != "20190902" ]
do
   curdate=`date -d"${startdate} 1 days" +"%Y%m%d"|awk '{printf($1)}'`
   startdate="$curdate"

   if [ -f /home/stock/stockdata/${curdate}/stockdata.txt ];then
     lclsize=`ls -ls /home/stock/stockdata/${curdate}"|grep stockdata.txt|awk -F ' ' '{printf $6}'`
     rmtsize=`ssh stock@47.89.240.84 "ls -ls /home/stock/stockdata/${curdate}"|grep stockdata.txt|awk -F ' ' '{printf $6}'`
     if [ $lclsize -ne $rmtsize ];then
         dataop $curdate
     fi
   else
     dataop $curdate
   fi

done
