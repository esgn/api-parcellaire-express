#! /bin/bash

# update script files in /tmp
/bin/bash ./update.sh

python3 -c $'import time\nwhile True:\n     time.sleep(3600)'

return 0