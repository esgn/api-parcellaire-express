#! /bin/bash

# ensures the volume contains the latest version of the scripts
cp -v /scripts/* /tmp/.

python3 -c $'import time\nwhile True:\n     time.sleep(3600)'

return 0