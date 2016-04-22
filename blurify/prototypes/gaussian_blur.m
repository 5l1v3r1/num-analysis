function [m] = gaussian_blur(d, stddev, spread)
  m = eye(d^2, d^2);
  divisor = computeTotalWeight(stddev, spread);
  for i = 1:d
    for j = 1:d
      sum = 0;
      rowIdx = pointIdx(d, i, j);
      for x = i-spread:i+spread
        for y = j-spread:j+spread
          weight = blurCoeff(x-i, y-j, stddev, spread);
          weight /= divisor;
          m(rowIdx, pointIdx(d, x, y)) = weight;
        end
      end
    end
  end
end

function [idx] = pointIdx(d, x, y)
  x = wrapIdx(d, x);
  y = wrapIdx(d, y);
  idx = x + (y-1)*d;
end

function [idx] = wrapIdx(d, x)
  idx = x;
  if (idx < 1)
    idx = d + idx;
  elseif (idx > d)
    idx = idx - d;
  end
end

function [s] = computeTotalWeight(stddev, spread)
  s = 0;
  for x = -spread:spread
    for y = -spread:spread
      s += blurCoeff(x, y, stddev, spread);
    end
  end
end

function [b] = blurCoeff(x, y, stddev, spread)
  sigSq = stddev^2;
  b = exp(-(x^2 + y^2)/(2*sigSq));
end
