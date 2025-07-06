
export function getDefaultProps(type: string): Record<string, any> {
	switch (type) {
		case 'heading':
			return {
				title: 'Heading Title',
				description: 'Optional subtext'
			};

		case 'image':
			return {
				mediaId: '',
				alt: '',
				size: 'medium',
				alignment: 'center',
				objectFit: 'cover',
				objectPosition: 'center'
			};


		case 'gallery':
			return {
				mediaIds: [],
				layout: 'grid',
				columns: 3,
				showCaptions: true,
				autoplay: false,
				gap: 16
			};
		case 'richtext':
			return {
				html: '<p>Start writing...</p>'
			};
		default:
			return {};
	}
}
