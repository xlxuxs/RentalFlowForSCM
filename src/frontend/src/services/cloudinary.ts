// Cloudinary Image Upload Service

const CLOUDINARY_CLOUD_NAME = 'ds6ednoms';
const CLOUDINARY_UPLOAD_PRESET = 'rentalflow';

interface CloudinaryUploadResponse {
    secure_url: string;
    public_id: string;
    width: number;
    height: number;
    format: string;
}

/**
 * Upload an image to Cloudinary
 * @param file - The image file to upload
 * @returns The secure URL of the uploaded image
 */
export async function uploadImage(file: File): Promise<string> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('upload_preset', CLOUDINARY_UPLOAD_PRESET);

    const response = await fetch(
        `https://api.cloudinary.com/v1_1/${CLOUDINARY_CLOUD_NAME}/image/upload`,
        {
            method: 'POST',
            body: formData,
        }
    );

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error('Cloudinary upload error:', errorData);
        throw new Error(errorData.error?.message || 'Failed to upload image');
    }

    const data: CloudinaryUploadResponse = await response.json();
    return data.secure_url;
}

/**
 * Upload multiple images to Cloudinary
 * @param files - Array of image files to upload
 * @returns Array of secure URLs
 */
export async function uploadImages(files: File[]): Promise<string[]> {
    const uploadPromises = files.map(file => uploadImage(file));
    return Promise.all(uploadPromises);
}

/**
 * Get an optimized Cloudinary URL with transformations
 * @param url - Original Cloudinary URL
 * @param options - Transformation options
 * @returns Optimized URL
 */
export function getOptimizedUrl(
    url: string,
    options: { width?: number; height?: number; quality?: number } = {}
): string {
    if (!url || !url.includes('cloudinary')) {
        return url;
    }

    const { width = 800, height, quality = 80 } = options;

    // Build transformation string
    const transforms = [`w_${width}`, `q_${quality}`, 'f_auto', 'c_fill'];
    if (height) {
        transforms.push(`h_${height}`);
    }

    // Insert transformations into URL
    return url.replace('/upload/', `/upload/${transforms.join(',')}/`);
}

/**
 * Get a thumbnail URL from Cloudinary
 */
export function getThumbnailUrl(url: string): string {
    return getOptimizedUrl(url, { width: 300, height: 200, quality: 70 });
}

// Alias for backward compatibility and Profile page
export { uploadImage as uploadToCloudinary };
