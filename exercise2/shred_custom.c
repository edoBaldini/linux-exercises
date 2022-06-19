#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

/* shred_custom.c aims at removing the specified file.
 * It overwrite the contents of the file three times:
 * - first time with 0
 * - second time with 1
 * - third time with a random value
 * After these three steps it removes the file.
 */


long min(long a, long b)
{
	if (a < b) 
		return a;
	return b;
}
int overwrite(int target, char* fn){
	int fd;
	char buf[4096] = {0};
	struct stat st;
	long index;
	long size_left;
	ssize_t written;
	fd=open(fn, O_WRONLY);
	if (fd < 0){	
		fprintf(stderr, "Provide a valid file name\n");
		exit(1);
	}
	switch (target)
	{
		case 1:
			for (int j=0; j < 4096; j++) {
				buf[j] = ~buf[j];
			}
			break;
		case -1:
			int frandom = open("/dev/random", O_RDONLY);
			read(frandom, buf, 4095);
			close(frandom);
			break;
	}
	fstat(fd, &st);
	for (index = 0; index < st.st_size; index += written)
        {
                size_left = min(st.st_size - index, 4096);
                written = write(fd, buf, size_left);
                if (written == -1){
                        fprintf(stderr, "Error occur: %s\n", strerror(errno));
                        exit(1);
                }
        }
	fsync(fd);
	close(fd);
	return 0;
}

int main(int argc, char *argv[]){
	char *fn;
	char buf[4096];
	int removed;
	if (argc == 2)
		fn = argv[1];
	else{
		printf("Provide a file name\n");
		exit(1);
	}
	overwrite(0, fn);
	overwrite(1, fn);
	overwrite(-1, fn);	
	
	removed = remove(fn);
	if (removed == 0) 
		printf("Removed\n");
	else
	{
		fprintf(stderr, "Error occur: %s\n", strerror(errno));
		return 1;
	}


	return 0;
}
