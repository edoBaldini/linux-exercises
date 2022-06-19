#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

long min(long a, long b)
{
	if (a < b) 
		return a;
	return b;
}
int overwrite(int target, char* fn, char *buf){
	int fd;
	struct stat st;
	long index;
	long size_left;
	ssize_t written;
	fd=open(fn, O_WRONLY);
	if (fd < 0){	
		fprintf(stderr, "Provide a valid file name\n");
		exit(1);
	}
	if (target == -1){
		int frandom = open("/dev/random", O_RDONLY);
		read(fd, buf, 4096);
		close(frandom);
	}
	else{
		memset(buf, target, 4096);
	}
	fstat(fd, &st);
	printf("%s\n", buf);
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
	if (argc == 2)
		fn = argv[1];
	else{
		printf("Provide a file name\n");
		exit(1);
	}
	overwrite(0, fn, buf);
	overwrite(1, fn, buf);
	overwrite(-1, fn, buf);	
	
	//remove_it;
	return 0;
}
