

int	printf(char *str)
{
	write(1, str, strlen(str));
}

int	*malloc(int length)
{
	t_page *list;
	if (length > 4096)
	{
		while (length)
		{
			void *page = list->get_page();
			list->page = mmap(0, 4096, PROT_READ | PROT_WRITE,
				MAP_ANON | MAP_PRIVATE, -1, 0);
			list = list->next;
			length -= 4096
		}
	}
	else 
	{
		list->page = mmap(0, length, PROT_READ | PROT_WRITE,
			MAP_ANON | MAP_PRIVATE, -1, 0);
	}
	return list->page->get_page();
}

int	write(int file, char *buffer, int length) 
{
	__asm__("mov eax, 4\n
		mov edi, %0\n
		mov esi, %1\n
		mov edx, %2\n
		syscall",
		: (file) (buffer) (length));
}

int	strlen(char *str)
{
	int i = 0;
	while (str[i]) {
		i++;
	}
	return (i);
}

char 	*buffer(int length) {
	value = (char *)malloc(length);
	if (value != NULL)
		return value;
	printf("Error occured while making buffer");
	return (0);
}
